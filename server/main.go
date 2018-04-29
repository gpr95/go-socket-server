package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"golang.org/x/crypto/blowfish"
	"encoding/hex"
)

/*
Creates socket and accept all connections.
Run handler function every time client approaches.
 */
func SocketServer(port int) {
	listener, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	defer listener.Close()
	if err != nil {
		log.Fatalf("Socket listener port %d failed,%s", port, err)
		os.Exit(1)
	}
	log.Printf("Begin listener port: %d", port)

	// Infinite loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	defer conn.Close()

	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
	)
	key := []byte("my key7")
	cipher, _ := blowfish.NewCipher(key)

	var decrypt [8]byte

ILOOP:
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			break ILOOP
		case nil:
			cipher.Decrypt(decrypt[0:], buf)
			log.Println(hex.Dump(decrypt[0:8]))
			if isTransportOver(data) {
				break ILOOP
			}

		default:
			log.Fatalf("Receive data failed:%s", err)
			return
		}
	}

}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, "\r\n\r\n")
	return
}

func main() {
	port := 443
	SocketServer(port)
}
