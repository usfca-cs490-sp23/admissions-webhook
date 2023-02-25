package keygen

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
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

	ca = x509.Certificate{
		SerialNumber: genSerialNum(),
		Subject: pkix.Name{
			Organization: []string{"usfca-cs490"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
)

// A cabundle is just the ca.crt (the CA's certificate) file in base64 format
func genCABundle() []byte {
	// NOTE: This function is based off of: Velotio's CAbundle tutorial (https://www.velotio.com/engineering-blog/managing-tls-certificate-for-kubernetes-admission-webhook)
	// CA private key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Make the CA cert (stored into caBytes)
	caBytes, err := x509.CreateCertificate(rand.Reader, &template, &ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Encode the ca cert to .pem format
	var caPEM = new(bytes.Buffer)
	_ = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	// get the result of encoding the CA cert into .pem (which is now stored in a byte slice)
	caRaw := caPEM.Bytes()
	return caRaw
}

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
func CreatePEMs() ([]byte, []byte, []byte) {
	// generate the CA and CAbundle
	caBundle := genCABundle()

	// Generate a private key
	privateKey := genPrivateKey()
	// Create the certificate in x509 format from privateKey
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &ca, &privateKey.PublicKey, privateKey)
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

	return pemCert, pemKey, caBundle
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

	}
}
