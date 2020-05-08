package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/krunalbhatt9/oniongeneration/onions"
)

const (
	//StopCharacter while reading it should stop at this
	StopCharacter = "\r\n\r\n"
)
const (
	message = "Ping"
)

var (
	selectedRouter onions.Router
	gcm            cipher.AEAD
	nonceSize      int
)

func initializeGCM() {

	key, error := hex.DecodeString("c8e63ff24118dee4dfdf5d865a913088a846e381b8b03ffb723df41f2b1e970f")
	if error != nil {
		fmt.Printf("Error reading key: %s\n", error.Error())
		os.Exit(1)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("Error reading key: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Key: %s\n", hex.EncodeToString(key))

	gcm, err = cipher.NewGCM(block)
	if err != nil {
		fmt.Printf("Error initializing AEAD: %s\n", err.Error())
		os.Exit(1)
	}

	nonceSize = gcm.NonceSize()
}

//SocketClient sends the message
func SocketClient(ip string, port int, ciphertext []byte) {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	a := []byte("Hell")
	b := []byte("127.0.0.1:3334")

	p := &onions.Message{onions.Encrypt(a, gcm, nonceSize), onions.Encrypt(b, gcm, nonceSize)}
	encoder.Encode(p)

	//conn.Write([]byte(ciphertext))
	//conn.Write([]byte(StopCharacter))
	//log.Printf("Send: %s", message)

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])

}

func main() {
	fmt.Println(len(os.Args))
	initializeGCM()
	b := []byte("hello")
	ciphertext := onions.Encrypt(b, gcm, nonceSize)
	var (
		ip   = "127.0.0.1"
		port = 3333
	)

	SocketClient(ip, port, ciphertext)

}
