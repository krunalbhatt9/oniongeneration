package onions

import (
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

//Routers an array of routers
type Routers struct {
	Routers []Router `json:"routers"`
}

//Router struct
type Router struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
	Key  string `json:"key"`
}

//RandBytes gens random byte
func RandBytes(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

//Encrypt s the file
func Encrypt(plaintext []byte, gcm cipher.AEAD, nonceSize int) (ciphertext []byte) {
	nonce := RandBytes(nonceSize)
	c := gcm.Seal(nil, nonce, plaintext, nil)
	return append(nonce, c...)
}

//Decrypt s the string
func Decrypt(ciphertext []byte, gcm cipher.AEAD, nonceSize int) (plaintext []byte, err error) {
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("Ciphertext too short")
	}
	nonce := ciphertext[0:nonceSize]
	msg := ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, msg, nil)
}

//Message structure
type Message struct {
	A []string
}
