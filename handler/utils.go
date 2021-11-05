package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"
	"log"
)

var cryptoKey = []byte("ABCDEFGHIJKLMNOP")

// =================== CFB ======================
func aesEncryptCFB(origData []byte, key []byte) (encrypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted
}
func aesDecryptCFB(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}

func idEncode(id string) string {
	idByte := []byte(id)
	encrypted := aesEncryptCFB(idByte, cryptoKey)
	return base32.StdEncoding.EncodeToString(encrypted)
}

func idDecode(code string) string {
	codeByte, err := base32.StdEncoding.DecodeString(code)
	log.Println(codeByte)
	if err != nil {
		//TODO return err
		fmt.Println(err)
		return ""
	}
	encrypted := aesDecryptCFB(codeByte, cryptoKey)
	return string(encrypted)
}
