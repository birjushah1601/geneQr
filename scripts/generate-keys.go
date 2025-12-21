package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("🔑 Generating RSA Keys for JWT...")
	
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Failed to generate private key:", err)
	}
	
	// Create keys directory if it doesn't exist
	if err := os.MkdirAll("keys", 0755); err != nil {
		log.Fatal("Failed to create keys directory:", err)
	}
	
	// Save private key
	privateKeyFile, err := os.Create("keys/jwt-private.pem")
	if err != nil {
		log.Fatal("Failed to create private key file:", err)
	}
	defer privateKeyFile.Close()
	
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		log.Fatal("Failed to write private key:", err)
	}
	
	fmt.Println("✅ Private key generated: keys/jwt-private.pem")
	
	// Save public key
	publicKeyFile, err := os.Create("keys/jwt-public.pem")
	if err != nil {
		log.Fatal("Failed to create public key file:", err)
	}
	defer publicKeyFile.Close()
	
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatal("Failed to marshal public key:", err)
	}
	
	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	
	if err := pem.Encode(publicKeyFile, publicKeyPEM); err != nil {
		log.Fatal("Failed to write public key:", err)
	}
	
	fmt.Println("✅ Public key generated: keys/jwt-public.pem")
	fmt.Println("\n🎉 JWT keys generated successfully!")
	fmt.Println("\n⚠️  IMPORTANT: Keep keys/jwt-private.pem secure and never commit to git")
}
