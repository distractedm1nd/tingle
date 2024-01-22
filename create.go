package main

import (
	"crypto/ed25519"
)

func CreateCmd() error {
	return nil
}

func Create() (string, string) {
	// Generate a new key pair using ed25519
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		// Handle the error according to your error handling policy
		// For example, you could log the error and return empty strings
		// log.Printf("Error generating keys: %v", err)
		return "", ""
	}
	// Convert the keys to strings to return
	encryptionKey := string(pubKey)
	decryptionKey := string(privKey)
	return encryptionKey, decryptionKey
}
