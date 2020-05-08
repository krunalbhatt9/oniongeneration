package onions

import (
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

var (
	gcm            cipher.AEAD
	nonceSize      int
	selectedRouter Router
)

const (
	StopCharacter = "\r\n\r\n"
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

func randBytes(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

func encrypt(plaintext []byte) (ciphertext []byte) {
	nonce := randBytes(nonceSize)
	c := gcm.Seal(nil, nonce, plaintext, nil)
	return append(nonce, c...)
}

func decrypt(ciphertext []byte) (plaintext []byte, err error) {
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("Ciphertext too short.")
	}
	nonce := ciphertext[0:nonceSize]
	msg := ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, msg, nil)
}
