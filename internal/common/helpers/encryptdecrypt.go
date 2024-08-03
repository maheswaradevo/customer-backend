package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

const (
	KeyEncrypDecryptData = "DataEncryptDecryptKeyIop"
)

func DecryptData(plainText string, key string) string {
	var err error
	if len(plainText) == 0 {
		return plainText
	}
	text, err := base64.StdEncoding.DecodeString(plainText)
	if err != nil {
		return ""
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err.Error()
	}
	if len(text) < aes.BlockSize {
		message := "size block invalid"
		return message
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	return string(text)
}

func EncryptData(plainText string, key string) string {
	var err error
	if len(plainText) == 0 {
		return plainText
	}
	text := []byte(plainText)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err.Error()
	}
	cipherText := make([]byte, aes.BlockSize+len(text))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return err.Error()
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], text)
	return base64.StdEncoding.EncodeToString(cipherText)
}
