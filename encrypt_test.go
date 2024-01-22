package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	key := "examplekey123456"
	msg := "hello world"
	encryptedText, err := Encrypt(msg, key)
	require.NoError(t, err)

	decryptedText, err := Decrypt(encryptedText, key)
	require.NoError(t, err)
	require.Equal(t, msg, decryptedText)
}
