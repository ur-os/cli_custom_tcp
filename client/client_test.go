package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"time"

	"testing"
)

var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces).")
var port = flag.Int("port", 8000, "The port to listen on; default is 8000.")

func TestUnit(t *testing.T) {
	go func() {
		flag.Parse()
		src := *addr + ":" + strconv.Itoa(*port)
		listener, err1 := net.Listen("tcp", src)
		if err1 != nil {
			fmt.Printf("Some connection error: %s\n", err1)
		}
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Some connection error: %s\n", err)
			}
			scanner := bufio.NewScanner(conn)

			for {
				ok := scanner.Scan()

				if !ok {
					break
				}

				message := scanner.Text()

				if message == "\u0000\u0000\u0000\u0002\u0000\u0000\u0000\u001A\u0000\u0000\u0000\u0001\u0000\u0000\u0000\u0002\u0000\u0000\u0000\vabracadabra\u0000\u0000\u0000\u0003xxx" {
					packet := buildPacketErrorResponse(
						1,
						2,
						1,
						[]byte("token not found"),
					)
					conn.Write(packet)
				}
			}
		}

	}()

	time.Sleep(3 * time.Second)
	out, err := exec.Command("./client", "loacalhost", "8000", "abracadabra", "xxx").Output()
	if err != nil {
		log.Fatal(err)
	}
	if string(out) != "if err != nil {\n    log.Fatal(err)\n}" {
		t.Errorf("Expect %s, got %s",
			"if err != nil {\n    log.Fatal(err)\n}",
			out)
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
