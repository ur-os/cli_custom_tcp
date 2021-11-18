/*
A very simple TCP server written in Go.

This is a toy project that I used to learn the fundamentals of writing
Go code and doing some really basic network stuff.

Maybe it will be fun for you to read. It's not meant to be
particularly idiomatic, or well-written for that matter.
*/
package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces).")
var port = flag.Int("port", 8000, "The port to listen on; default is 8000.")

const MinInt8, MaxInt8 = -(-128), 127
const MinInt32, MaxInt32 = -(-2147483648), 2147483647
const MinInt64, MaxInt64 = -(-9223372036854775808), 9223372036854775807

func main() {
	flag.Parse()

	fmt.Println("Starting server...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Listening on %s.\n", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)

	scanner := bufio.NewScanner(conn)

	for {
		ok := scanner.Scan()

		if !ok {
			break
		}

		handleMessage(scanner.Text(), conn)
	}

	fmt.Println("Client at " + remoteAddr + " disconnected.")
}

func handleMessage(message string, conn net.Conn) {
	fmt.Println("> " + message)

	if len(message) > 0 {
		switch {
		case message == "/time":
			resp := "It is " + time.Now().String() + "\n"
			fmt.Print("< " + resp)
			conn.Write([]byte(resp))

		case message == "/quit":
			fmt.Println("Quitting.")
			conn.Write([]byte("I'm shutting down now.\n"))
			fmt.Println("< " + "%quit%")
			conn.Write([]byte("%quit%\n"))
			os.Exit(0)
		case message == "abracadabra xxx":
			fmt.Println("CUBE_OAUTH2_ERR_BAD_SCOPE")
			conn.Write([]byte("error: CUBE_OAUTH2_ERR_BAD_SCOPE\n message: bad scope\n"))
		case message == "1":
			fmt.Println("TEST_BINARY_STREAM")

			var packet []byte
			//bsInt8 := make([]byte, 1)
			//bsInt32 := make([]byte, 4)
			//bsInt64 := make([]byte, 8)

			packet = buildPacket(
				-15,
				0x00000002,
				0x00000001,
				[]byte("fkjas;fosdjfofaso;fdjsofasdfisdjfoidsfjoa"),
				0x00000001,
				[]byte("ur_0s"),
				0x00200001,
				0x0000000000000001,
			)

			size, _ := conn.Write(packet)
			fmt.Println(size)
			fmt.Println("TEST_BINARY_STREAM_END")
		default:
			conn.Write([]byte("Unrecognized command.\n"))
		}
	}
}

func buildPacket(
	svcId int32,
	requestId int32,
	returnCode int32,
	clientId []byte,
	clientType int32,
	username []byte,
	expiresIn int32,
	userId int64,
) []byte {

	//  TODO: to process exceptions

	//bsInt8 := make([]byte, 1)
	bsInt32 := make([]byte, 4)
	bsInt64 := make([]byte, 8)

	var packet []byte
	var body []byte
	//var header []byte

	//
	// body
	//
	binary.BigEndian.PutUint32(bsInt32, uint32(returnCode))
	for _, character := range bsInt32 {
		body = append(body, character)
	}

	binary.BigEndian.PutUint32(bsInt32, uint32(len(clientId)))
	for _, character := range bsInt32 {
		body = append(body, character)
	}
	for _, oneByte := range clientId {
		body = append(body, oneByte)
	}

	binary.BigEndian.PutUint32(bsInt32, uint32(clientType))
	for _, character := range bsInt32 {
		body = append(body, character)
	}

	binary.BigEndian.PutUint32(bsInt32, uint32(len(username)))
	for _, character := range bsInt32 {
		body = append(body, character)
	}
	for _, oneByte := range username {
		body = append(body, oneByte)
	}

	binary.BigEndian.PutUint32(bsInt32, uint32(expiresIn))
	for _, character := range bsInt32 {
		body = append(body, character)
	}

	binary.BigEndian.PutUint64(bsInt64, uint64(userId))
	for _, character := range bsInt64 {
		body = append(body, character)
	}

	//
	// header
	//
	binary.BigEndian.PutUint32(bsInt32, uint32(svcId))
	for _, character := range bsInt32 {
		packet = append(packet, character)
	}
	fmt.Println(svcId)

	//  bodyLength
	binary.BigEndian.PutUint32(bsInt32, uint32(len(body)))
	for _, character := range bsInt32 {
		packet = append(packet, character)
	}

	binary.BigEndian.PutUint32(bsInt32, uint32(requestId))
	for _, character := range bsInt32 {
		packet = append(packet, character)
	}

	//for _, character := range body {
	//	packet = append(packet, character)
	//}

	packet = append(packet, body...)

	fmt.Println(packet)

	return packet
}
