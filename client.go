package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	routers   onions.Routers
	gcm       cipher.AEAD
	nonceSize int
)

func initializeGCM(ptr *int) {

	var selectedRouter = routers.Routers[*ptr]
	key, error := hex.DecodeString(selectedRouter.Key)
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
func SendOnion() {
	ip := "127.0.0.1"
	port := 3333

	a := []string{"Hell", "127.0.0.1:3335", "127.0.0.1:3334"}
	for i := range a {

		idx := len(a) - i - 1
		log.Printf("i %d", i)
		log.Printf("idx %d", idx)
		initializeGCM(&idx)
		for i2 := 0; i2 <= i; i2++ {
			a[i2] = string(onions.Encrypt([]byte(a[i2]), gcm, nonceSize))
			log.Printf("Encrypting")
		}

	}
	//b := routers.Routers[i].IP + ":" + strconv.Itoa(routers.Routers[i].Port)
	//log.Printf("Sending to ", b)

	p := &onions.Message{a}
	log.Printf("Send: ", p)

	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", string(addr))

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer conn.Close()

	encoder := gob.NewEncoder(conn)

	// 	fmt.Println("User Type: " + routers.Routers[i].IP)
	// 	fmt.Println("User Age: " + strconv.Itoa(routers.Routers[i].Port))
	// 	fmt.Println("User Name: " + routers.Routers[i].Name)

	encoder.Encode(p)

	//conn.Write([]byte(ciphertext))
	//conn.Write([]byte(StopCharacter))
	log.Printf("Send: ", p, conn)

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])

}

func main() {
	jsonFile, err := os.Open("properties.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened properties.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &routers)
	// for i := 0; i < len(routers.Routers); i++ {
	// 	fmt.Println("User Type: " + routers.Routers[i].IP)
	// 	fmt.Println("User Age: " + strconv.Itoa(routers.Routers[i].Port))
	// 	fmt.Println("User Name: " + routers.Routers[i].Name)
	// }
	//selectedRouter = routers.Routers[*ptr]

	// initializeGCM()
	// b := []byte("hello")
	// ciphertext := onions.Encrypt(b, gcm, nonceSize)
	// var (
	// 	ip   = "127.0.0.1"
	// 	port = 3333
	// )

	SendOnion()

}
