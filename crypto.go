// Package goshared ...
package goshared

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func AES256Encrypt(stringToEncrypt string, keyString string) (encryptedString string) {
	//Since the key is in string, we need to convert decode it to bytes
	key, err := hex.DecodeString(keyString)
	if err != nil {
		panic(err)
	}
	if len(key) != 32 { // AES256 key length is 32
		panic(fmt.Sprintf("key length should be 32 in byte, your key string has %d byte", len(key)))
	}
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func AES256Decrypt(encryptedString string, keyString string) (decryptedString string) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		panic(err.Error())
	}
	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		panic(err.Error())
	}

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%s", plaintext)
}
