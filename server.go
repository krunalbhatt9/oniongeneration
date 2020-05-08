package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/krunalbhatt9/oniongeneration/onions"
)

var (
	selectedRouter onions.Router
	gcm            cipher.AEAD
	nonceSize      int
)

const (
	//StopCharacter while reading it should stop at this
	StopCharacter = "\r\n\r\n"
)

func initializeGCM() {

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
func SocketClient(IP string, port int, message string) {
	addr := strings.Join([]string{IP, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
		log.Printf("Ip address could not be resolved. The message has reached the last node in the chain")
		//os.Exit(1)
	}

	defer conn.Close()

	conn.Write([]byte(message))
	conn.Write([]byte(StopCharacter))
	log.Printf("Send: %s", message)

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])

}

//ReadandSendMessage reads and sends the message
func ReadandSendMessage(message []byte) {
	plaintext, err := onions.decrypt(message, gcm, nonceSize)
	SocketClient("127.0.0.1", 3334, plaintext)
}

//SocketServer sends the message
func SocketServer(IP string, port int) {

	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
		os.Exit(1)
	}

	defer listen.Close()

	log.Printf("Begin listen port: %d", port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go handler(conn)
	}

}

// func handleConnection(conn net.Conn) {
// 	dec := gob.NewDecoder(conn)
// 	p := &Message{}
// 	dec.Decode(p)
// 	fmt.Printf("Received : %+v", p.IP)
// 	ReadandSendMessage(p.IP, p.Port, p.Data)
// 	conn.Close()
// }

func handler(conn net.Conn) {

	defer conn.Close()

	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
		w   = bufio.NewWriter(conn)
	)

ILOOP:
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			break ILOOP
		case nil:
			log.Println("Receive:", data)
			if isTransportOver(data) {
				ReadandSendMessage(data)
				break ILOOP
			}

		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}

	}
	w.Write([]byte("OK"))
	w.Flush()
	log.Printf("Send: %s", "OK")

}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\r\n\r\n")
	return
}

func main() {
	ptr := flag.Int("router", 0, "an int")
	flag.Parse()
	fmt.Println("router:", *ptr)
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
	selectedRouter = routers.Routers[*ptr]
	initializeGCM()
	SocketServer(selectedRouter.IP, selectedRouter.Port)
}
