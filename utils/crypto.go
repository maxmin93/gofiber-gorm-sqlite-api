package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

var (
	initialVector = "1234567890123456"
	passphrase    = "Impassphrasegood"
)

// Crypto is main function
func Crypto() {
	var plainText = "hello world"

	encryptedData := AESEncrypt(plainText, []byte(passphrase))
	encryptedString := base64.StdEncoding.EncodeToString(encryptedData)
	fmt.Println(encryptedString)

	encryptedData, _ = base64.StdEncoding.DecodeString(encryptedString)
	decryptedText := AESDecrypt(encryptedData, []byte(passphrase))
	fmt.Println(string(decryptedText))
}

// AESEncrypt encrypts a string to a base64 encoded string
func AESEncrypt(src string, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCEncrypter(block, []byte(initialVector))
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)

	return crypted
}

// AESDecrypt decrypts a base64 encoded string to a string
func AESDecrypt(crypt []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if len(crypt) == 0 {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCDecrypter(block, []byte(initialVector))
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)

	return PKCS5Trimming(decrypted)
}

// PKCS5Padding pads the data to the block size
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5Trimming removes the padding from the data
func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
