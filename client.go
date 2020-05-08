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

	"github.com/krunalbhatt9/oniongeneration/onions"
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
		fmt.Printf("Client: Error reading key: %s\n", error.Error())
		os.Exit(1)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("Client: Error reading key: %s\n", err.Error())
		os.Exit(1)
	}

	//fmt.Printf("Key: %s\n", hex.EncodeToString(key))

	gcm, err = cipher.NewGCM(block)
	if err != nil {
		fmt.Printf("Client: Error initializing AEAD: %s\n", err.Error())
		os.Exit(1)
	}

	nonceSize = gcm.NonceSize()
}

//SendOnion sends the clients encypted onion
func SendOnion(a []string, addr string) {

	for i := range a {

		idx := len(a) - i - 1
		//log.Printf("i %d", i)
		//log.Printf("idx %d", idx)
		log.Printf("Client: Encrypting the onion with keys of %s", routers.Routers[idx].IP+":"+strconv.Itoa(routers.Routers[idx].Port))
		initializeGCM(&idx)
		for i2 := 0; i2 <= i; i2++ {
			a[i2] = string(onions.Encrypt([]byte(a[i2]), gcm, nonceSize))
			//log.Printf("Encrypting")
		}

	}
	//b := routers.Routers[i].IP + ":" + strconv.Itoa(routers.Routers[i].Port)
	//log.Printf("Sending to ", b)

	p := &onions.Message{a}
	log.Printf("Client: Sending the onion too %s", addr)

	//addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
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
	//log.Printf("Send: ", p, conn)

	// buff := make([]byte, 1024)
	// n, _ := conn.Read(buff)
	// log.Printf("Receive: %s", buff[:n])

}

func main() {
	jsonFile, err := os.Open("properties.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Client: Successfully Opened properties.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	a := []string{"Hello this is onion routing!"}
	fmt.Println("Client: We are going to send the message:", a[0])
	json.Unmarshal(byteValue, &routers)
	for i := len(routers.Routers) - 1; i > 0; i-- {
		// fmt.Println("User Type: " + routers.Routers[i].IP)
		// fmt.Println("User Age: " + strconv.Itoa(routers.Routers[i].Port))
		// fmt.Println("User Name: " + routers.Routers[i].Name)
		a = append(a, routers.Routers[i].IP+":"+strconv.Itoa(routers.Routers[i].Port))
	}
	//selectedRouter = routers.Routers[*ptr]
	address := routers.Routers[0].IP + ":" + strconv.Itoa(routers.Routers[0].Port)
	// initializeGCM()
	// b := []byte("hello")
	// ciphertext := onions.Encrypt(b, gcm, nonceSize)
	// var (
	// 	ip   = "127.0.0.1"
	// 	port = 3333
	// )
	//a := []string{"Hell", "127.0.0.1:3335", "127.0.0.1:3334"}
	fmt.Println("Client: Listing the addresses of the nodes to use", a[1:], "\nEntry Node: ", address, "Destination Node: ", a[1])
	SendOnion(a, address)

}
