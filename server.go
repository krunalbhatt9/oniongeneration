package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/krunalbhatt9/oniongeneration/onions"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))
var (
	selectedRouter onions.Router
	gcm            cipher.AEAD
	nonceSize      int
)

func initializeGCM() {

	key, error := hex.DecodeString(selectedRouter.Key)
	if error != nil {
		fmt.Printf("Router %s: Error reading key: %s\n", selectedRouter.Name, error.Error())
		os.Exit(1)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("Router %s: Error reading key: %s\n", selectedRouter.Name, err.Error())
		os.Exit(1)
	}

	//fmt.Printf("Key: %s\n", hex.EncodeToString(key))

	gcm, err = cipher.NewGCM(block)
	if err != nil {
		fmt.Printf("Router %s: Error initializing AEAD: %s\n", selectedRouter.Name, err.Error())
		os.Exit(1)
	}

	nonceSize = gcm.NonceSize()
}

//SocketClient sends the message
func SocketClient(message []string, addr string) {
	//addr := strings.Join([]string{IP, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Router %s: Ip address could not be resolved. The message has reached the last node in the chain", selectedRouter.Name)
		log.Printf("Router %s: The message is :%s", selectedRouter.Name, addr)
		//log.Fatalln(err)
		//os.Exit(1)
		return
	}

	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	// a := []byte("Hell")
	// b := []byte("add")

	p := &onions.Message{message}
	encoder.Encode(p)

	//conn.Write([]byte(message))
	//conn.Write([]byte(StopCharacter))
	//log.Printf("Send: %s", message)

	// buff := make([]byte, 1024)
	// n, _ := conn.Read(buff)
	// log.Printf("Receive: %s", buff[:n])

}

//ReadandSendMessage reads and sends the message
func ReadandSendMessage(message []string) {
	log.Printf("Router %s: Decrypting the packet", selectedRouter.Name)
	length := len(message) - 1
	for i, s := range message {
		s, err := onions.Decrypt([]byte(s), gcm, nonceSize)
		message[i] = string(s)
		if err != nil {
			log.Printf("Router %s: Failed to decrypt.Message %s recived", selectedRouter.Name, message)
			return
			//os.Exit(1)
		}
	}

	addr := message[length]
	// token := make([]byte, 32)
	// rand.Read(token)
	// randomString := string(token)
	// message = append(randomString, message)
	log.Printf("Router %s: ReRouting the packet to %s", selectedRouter.Name, addr)
	message = message[:length]
	random := strconv.Itoa(rand.Int())
	encryptedRandom := []string{string(onions.Encrypt([]byte(random), gcm, nonceSize))}

	message = append(encryptedRandom, message...)
	SocketClient(message, addr)

}

//SocketServer sends the message
func SocketServer(IP string, port int) {

	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("Router %s: Socket listen port %d failed, %s", selectedRouter.Name, port, err)
		os.Exit(1)
	}

	defer listen.Close()

	log.Printf("Router %s : Begin listen port: %d", selectedRouter.Name, port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	p := &onions.Message{}
	dec.Decode(p)
	//fmt.Println("Received : %d", len(p.A))
	ReadandSendMessage(p.A)
	conn.Close()
}

// func handler(conn net.Conn) {

// 	defer conn.Close()

// 	var (
// 		buf = make([]byte, 1024)
// 		r   = bufio.NewReader(conn)
// 		w   = bufio.NewWriter(conn)
// 	)

// ILOOP:
// 	for {
// 		n, err := r.Read(buf)
// 		data := string(buf[:n])

// 		switch err {
// 		case io.EOF:
// 			break ILOOP
// 		case nil:
// 			log.Println("Receive:", data)
// 			if isTransportOver(data) {
// 				b := []byte(data)
// 				ReadandSendMessage(b)
// 				break ILOOP
// 			}

// 		default:
// 			log.Fatalf("Receive data failed:%s", err)
// 			return
// 		}

// 	}
// 	w.Write([]byte("OK"))
// 	w.Flush()
// 	log.Printf("Send: %s", "OK")

// }

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
	//log.Printf("Router: Successfully Opened properties.json")
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
