package keygen

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

// NOTE: Based off of https://eli.thegreenplace.net/2021/go-https-servers-with-tls/

var (
	// General template for x509 Certificate with our desired settings
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
	// Test commit
)

/* Method to generate a private ecdsa private key */
func genPrivateKey() *ecdsa.PrivateKey {
	// Generate key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// Check for errors
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	return privateKey
}

/* Method to generate a random serial number */
func genSerialNum() *big.Int {
	// Maximum limit for serial num length = 128
	maxLenLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	// Generate a random serial number
	serialNumber, err := rand.Int(rand.Reader, maxLenLimit)
	// Check for errors
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	return serialNumber
}

/* Exported method to create cert.pem and key.pem files */
func CreatePem(certPath, keyPath string) ([]byte, []byte) {
	// Generate a private key
	privateKey := genPrivateKey()
	// Create the certificate in x509 format from privateKey
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	// Check for errors
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	// Encode x509 certificates to pem block
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	// Error check
	if pemCert == nil {
		log.Fatal("Failed to encode certificate to PEM")
	}

	// Generate x509 private key from privateKey
	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	// Error check
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	// Encode x509 private key to pem block
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes})
	// Error check
	if pemKey == nil {
		log.Fatal("Failed to encode key to PEM")
	}

	return pemCert, pemKey
}

/* Convert any file's contents to base64 using base64 utility */
func ConvertPEMToB64(data map[string][]byte) {
	// Loop through map
	for outfile, pemData := range data {
		// Encode PEM data from map
		data := make([]byte, base64.StdEncoding.EncodedLen(len(pemData)))
		base64.StdEncoding.Encode(data, pemData)

		// Error check
		if err := os.WriteFile(outfile, data, 0600); err != nil {
			log.Fatal(err)
		}

		log.Printf("wrote %s\n", outfile)
	}
}
