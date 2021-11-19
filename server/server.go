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
		case message == "\u0000\u0000\u0000\u0002\u0000\u0000\u0000\u001A\u0000\u0000\u0000\u0001\u0000\u0000\u0000\u0002\u0000\u0000\u0000\vabracadabra\u0000\u0000\u0000\u0003xxx":
			packet := buildPacketErrorResponse(
				1,
				2,
				1,
				[]byte("token not found"),
			)
			conn.Write(packet)
		case message == "\u0000\u0000\u0000\u0002\u0000\u0000\u0000\u001B\u0000\u0000\u0000\u0001\u0000\u0000\u0000\u0002\u0000\u0000\u0000\vabracadabra\u0000\u0000\u0000\u0004test":
			packet := buildPacketOkResponse(
				-15,
				0x00000002,
				0,
				[]byte("fkjas;fosdjfofaso;fdjsofasdfisdjfoidsfjoa"),
				0x00000001,
				[]byte("ur_0s"),
				0x00200001,
				0x0000000000000001,
			)

			conn.Write(packet)
		default:
			packet := buildPacketErrorResponse(
				1,
				2,
				1337,
				[]byte("Sorry im just a mock-server. Sun is white, river is blue; mocked functions send const's you"),
			)
			conn.Write(packet)
		}
	}
}

func buildPacketOkResponse(
	svcId int32,
	requestId int32,

	returnCode int32,
	clientId []byte,
	clientType int32,
	username []byte,
	expiresIn int32,
	userId int64,
) []byte {

	bsInt32 := make([]byte, 4)
	bsInt64 := make([]byte, 8)

	var packet []byte
	var body []byte

	//
	//  body filling
	//
	binary.BigEndian.PutUint32(bsInt32, uint32(returnCode))
	body = append(body, bsInt32...)

	binary.BigEndian.PutUint32(bsInt32, uint32(len(clientId)))
	body = append(body, bsInt32...)
	body = append(body, clientId...)

	binary.BigEndian.PutUint32(bsInt32, uint32(clientType))
	body = append(body, bsInt32...)

	binary.BigEndian.PutUint32(bsInt32, uint32(len(username)))
	body = append(body, bsInt32...)
	body = append(body, username...)

	binary.BigEndian.PutUint32(bsInt32, uint32(expiresIn))
	body = append(body, bsInt32...)

	binary.BigEndian.PutUint64(bsInt64, uint64(userId))
	body = append(body, bsInt64...)

	//
	//  header filling
	//
	binary.BigEndian.PutUint32(bsInt32, uint32(svcId))
	packet = append(packet, bsInt32...)

	binary.BigEndian.PutUint32(bsInt32, uint32(len(body))) // bodyLength
	packet = append(packet, bsInt32...)

	binary.BigEndian.PutUint32(bsInt32, uint32(requestId))
	packet = append(packet, bsInt32...)

	//
	//  packet = header + body
	//
	packet = append(packet, body...)

	return packet
}

func buildPacketErrorResponse(
	svcId int32,
	requestId int32,

	returnCode int32,
	errorString []byte,
) []byte {

	bsInt32 := make([]byte, 4)

	var packet []byte
	var body []byte

	//
	//  body filling
	//
	binary.BigEndian.PutUint32(bsInt32, uint32(returnCode))
	body = append(body, bsInt32...)

	binary.BigEndian.PutUint32(bsInt32, uint32(len(errorString)))
	body = append(body, bsInt32...)
	body = append(body, errorString...)

	//
	//  header filling
	//
	binary.BigEndian.PutUint32(bsInt32, uint32(svcId))
	packet = append(packet, bsInt32...)

	binary.BigEndian.PutUint32(bsInt32, uint32(len(body))) // bodyLength
	packet = append(packet, bsInt32...)

	binary.BigEndian.PutUint32(bsInt32, uint32(requestId))
	packet = append(packet, bsInt32...)

	//
	//  packet = header + body
	//
	packet = append(packet, body...)

	return packet
}

type server interface {
	Write(packet []byte) error
	Read(packet []byte) error
}
