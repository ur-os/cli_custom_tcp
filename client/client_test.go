package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"io/ioutil"
	"net"
	"testing"
)

var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces).")
var port = flag.Int("port", 8000, "The port to listen on; default is 8000.")

func TestUnit(t *testing.T) {

	packet := buildPacketErrorResponse(
		1,
		2,
		1,
		[]byte("token not found"),
	)

	go func() {
		conn, err := net.Dial("tcp", ":3000")
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		conn.Write(packet)
		//if _, err := fmt.Fprintf(conn, string(packet)); err != nil {
		//	t.Fatal(err)
		//}
	}()

	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			t.Fatal(err)
		}

		result := bytes.Compare(buf, packet)
		if result != 0 {
			t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", buf, packet)
		}
		return // Done
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
