package main
import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	message       = "Ping"
	StopCharacter = "\r\n\r\n"
)

//Message is a message to be sent
type Message struct {
	IP   string
	Port int
	Data string
}

//SocketClient sends the message
func SocketClient(ip string, port int) {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer conn.Close()
	encoder := gob.NewEncoder(conn)
	p := &Message{ip, port, message}
	encoder.Encode(p)
	//conn.Write([]byte("Ping"))
	//conn.Write([]byte(StopCharacter))
	//log.Printf("Send: %s", message)

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])

}

func main() {
	fmt.Println(len(os.Args))
	var (
		ip   = "127.0.0.1"
		port = 3333
	)

	SocketClient(ip, port)

}
