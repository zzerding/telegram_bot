package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"io"
)

var cryptoKey = []byte(BotToken[0:32])

// =================== CFB ======================
func aesEncryptCFB(origData []byte, key []byte) (encrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return
}
func aesDecryptCFB(encrypted []byte, key []byte) (decrypted []byte, err error) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		err = errors.New("ciphertext too short")
		return
	}
	iv := encrypted[:aes.BlockSize]
	decrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decrypted, decrypted)
	return
}

func idEncode(id string) (idEncode string, err error) {
	idByte := []byte(id)
	var encrypted []byte
	encrypted, err = aesEncryptCFB(idByte, cryptoKey)
	if err != nil {
		return
	}
	idEncode = base32.StdEncoding.EncodeToString(encrypted)
	return
}

func idDecode(code string) (id string, err error) {
	var codeByte, encrypted []byte
	codeByte, err = base32.StdEncoding.DecodeString(code)
	if err != nil {
		return
	}
	encrypted, err = aesDecryptCFB(codeByte, cryptoKey)
	if err != nil {
		return
	}
	id = string(encrypted)
	return
}
