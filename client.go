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

func initializeGCM(ptr int) {
	var routers onions.Routers
	selectedRouter = routers.Routers[ptr]
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

	//for i := 0; i < len(routers.Routers); i++ {
	// 	fmt.Println("User Type: " + routers.Routers[i].IP)
	// 	fmt.Println("User Age: " + strconv.Itoa(routers.Routers[i].Port))
	// 	fmt.Println("User Name: " + routers.Routers[i].Name)
	//}
	//addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", "127.0.0.1:3333")

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	a := []byte("Hell")
	b := []byte("127.0.0.1:3334")
	initializeGCM(1)
	c := onions.Encrypt(a, gcm, nonceSize)
	initializeGCM(0)
	p := &onions.Message{onions.Encrypt(c, gcm, nonceSize), onions.Encrypt(b, gcm, nonceSize)}
	encoder.Encode(p)

	//conn.Write([]byte(ciphertext))
	//conn.Write([]byte(StopCharacter))
	//log.Printf("Send: %s", message)

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

	var routers onions.Routers
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
