package keygen

import (
    "crypto/rand"
    "crypto/elliptic"
    "crypto/ecdsa"
    "crypto/x509"
    "crypto/x509/pkix"
    "time"
    "log"
    "math/big"
    "os"
    "encoding/pem"
)

// NOTE: Based off of https://eli.thegreenplace.net/2021/go-https-servers-with-tls/

var (
    template = x509.Certificate{
        SerialNumber: genSerialNum(),
        Subject: pkix.Name{
            Organization: []string{"usfca-cs490"},
        },
        DNSNames:  []string{"localhost"},
        NotBefore: time.Now(),
        NotAfter:  time.Now().Add(3 * time.Hour),

        KeyUsage:              x509.KeyUsageDigitalSignature,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
        BasicConstraintsValid: true,
    }
)

func genPrivateKey() (*ecdsa.PrivateKey) {
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        log.Fatalf("Failed to generate private key: %v", err)
    }

    return privateKey
}

func genSerialNum() (*big.Int) {
    serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
    serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
    if err != nil {
        log.Fatalf("Failed to generate serial number: %v", err)
    }

    return serialNumber
}

func CreatePem() {
    privateKey := genPrivateKey()
    dataBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
    if err != nil {
        log.Fatalf("Failed to create certificate: %v", err)
    }

    pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: dataBytes})
    if pemCert == nil {
        log.Fatal("Failed to encode certificate to PEM")
    }
    if err := os.WriteFile("cert.pem", pemCert, 0644); err != nil {
        log.Fatal(err)
    }
    log.Print("wrote cert.pem\n")

    privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
    if err != nil {
        log.Fatalf("Unable to marshal private key: %v", err)
    }
    pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
    if pemKey == nil {
        log.Fatal("Failed to encode key to PEM")
    }
    if err := os.WriteFile("key.pem", pemKey, 0600); err != nil {
        log.Fatal(err)
    }
    log.Print("wrote key.pem\n")
}
