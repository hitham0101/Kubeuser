package pkg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"
)

func GeneratePrivateKey(userName string) {
	// Change this to your desired file name
	fileName := userName + ".key"

	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Println("Error generating RSA private key:", err)
		return
	}

	// Encode private key to PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Create file to write the private key
	keyFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating key file:", err)
		return
	}
	defer keyFile.Close()

	// Write PEM data to the file
	err = pem.Encode(keyFile, privateKeyPEM)
	if err != nil {
		fmt.Println("Error writing PEM data:", err)
		return
	}

	fmt.Println("RSA private key generated successfully:", fileName)
}

func GenerateCSR(userName string) {
	// Change these to your desired file names
	keyFileName := userName + ".key"
	csrFileName := userName + ".csr"

	// Read RSA private key from file
	keyFile, err := os.Open(keyFileName)
	if err != nil {
		fmt.Println("Error opening key file:", err)
		return
	}
	defer keyFile.Close()

	keyBytes, err := io.ReadAll(keyFile)
	if err != nil {
		fmt.Println("Error reading key file:", err)
		return
	}

	// Decode PEM encoded private key
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		fmt.Println("Invalid PEM encoded key")
		return
	}

	// Parse RSA private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing private key:", err)
		return
	}

	// Create CSR template
	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   userName,
			Organization: []string{"system:demo"},
		},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	// Create CSR
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, privateKey)
	if err != nil {
		fmt.Println("Error creating CSR:", err)
		return
	}

	// Encode CSR to PEM format
	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDER,
	})

	// Write CSR to file
	csrFile, err := os.Create(csrFileName)
	if err != nil {
		fmt.Println("Error creating CSR file:", err)
		return
	}
	defer csrFile.Close()

	_, err = csrFile.Write(csrPEM)
	if err != nil {
		fmt.Println("Error writing CSR to file:", err)
		return
	}

	fmt.Println("CSR generated successfully:", csrFileName)
}

func GenerateCertificate(userName string) {

	csrFileName := userName + ".csr"
	crtFileName := userName + ".crt"
	caCertFileName := "ca.crt"
	caKeyFileName := "ca.key"

	// Read CSR from file
	csrBytes, err := os.ReadFile(csrFileName)
	if err != nil {
		fmt.Println("Error reading CSR file:", err)
		return
	}

	// Decode PEM encoded CSR
	block, _ := pem.Decode(csrBytes)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		fmt.Println("Invalid PEM encoded CSR")
		return
	}

	// Parse CSR
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing CSR:", err)
		return
	}

	// Read CA certificate from file
	caCertBytes, err := os.ReadFile(caCertFileName)
	if err != nil {
		fmt.Println("Error reading CA certificate file:", err)
		return
	}

	// Decode PEM encoded CA certificate
	block, _ = pem.Decode(caCertBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		fmt.Println("Invalid PEM encoded CA certificate")
		return
	}

	// Parse CA certificate
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing CA certificate:", err)
		return
	}

	// Read CA private key from file
	caKeyBytes, err := os.ReadFile(caKeyFileName)
	if err != nil {
		fmt.Println("Error reading CA private key file:", err)
		return
	}

	// Decode PEM encoded CA private key
	block, _ = pem.Decode(caKeyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		fmt.Println("Invalid PEM encoded CA private key")
		return
	}

	// Parse CA private key
	caKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing CA private key:", err)
		return
	}

	// Generate a new serial number
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		fmt.Println("Error generating serial number:", err)
		return
	}

	// Create certificate template
	template := &x509.Certificate{
		Subject:               csr.Subject,
		SerialNumber:          serialNumber,
		PublicKey:             csr.PublicKey,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, template, caCert, csr.PublicKey, caKey)
	if err != nil {
		fmt.Println("Error creating certificate:", err)
		return
	}

	// Encode certificate to PEM format
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	// Write certificate to file
	certFile, err := os.Create(crtFileName)
	if err != nil {
		fmt.Println("Error creating certificate file:", err)
		return
	}
	defer certFile.Close()

	_, err = certFile.Write(certPEM)
	if err != nil {
		fmt.Println("Error writing certificate to file:", err)
		return
	}

	fmt.Println("Certificate generated successfully:", crtFileName)
}
