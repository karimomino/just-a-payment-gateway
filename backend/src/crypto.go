package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func encryptPAN(plaintext []byte) ([]byte, []byte, error) {
	key := []byte(os.Getenv("PAN_CRYPTO_KEY"))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

func HashCardToBytes(pan string) []byte {
	secretKey := []byte(os.Getenv("PAN_CRYPTO_KEY"))
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(pan))
	return mac.Sum(nil)
}

func decryptPAN(encryptedPayload []byte) ([]byte, error) {
	key := []byte(os.Getenv("PAN_CRYPTO_KEY"))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(encryptedPayload) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := encryptedPayload[:nonceSize], encryptedPayload[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
