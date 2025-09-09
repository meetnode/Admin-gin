package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

var secret = []byte(os.Getenv("APP_SECRET"))

// Encrypt text with AES
func Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	plain := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(plain))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plain)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt text with AES
func Decrypt(cryptoText string) (string, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
