package main

import (
	"encoding/binary"
	"fmt"
	"os/exec"
	"testing"
)

func TestOkResponse(t *testing.T) {
	packet := buildPacketOkResponse(
		0x00000002,
		0x00000001,
		0x00000000,
		[]byte("UID_0001337"),
		0x00000001,
		[]byte("richard@mailer.ru.com"),
		0x00200001,
		0x0000000000000001,
	)

	server := exec.Command("/home/urick0s/go/src/testMail/dummy_server/dummy_server", string(packet))
	go server.Run()

	out, err := exec.Command("./client", "localhost", "8000", "abracadabra", "test").Output()
	if err != nil {
		t.Fatal(err)
	}

	if string(out) != "return_сode: 0\n"+
		"client_id: UID_0001337\n"+
		"client_type: 1\n"+
		"username: richard@mailer.ru.com\n"+
		"expires_in: 2097153\n"+
		"user_id: 1" {
		t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", out, packet)
	}
}

func TestTokenNotFound(t *testing.T) {
	packet := buildPacketErrorResponse(
		1,
		2,
		1,
		[]byte("token not found"),
	)

	server := exec.Command("./client", string(packet))
	go server.Run()

	out, err := exec.Command("./client", "localhost", "8000", "abracadabra", "xxx").Output()
	if err != nil {
		t.Fatal(err)
	}

	if string(out) != "error: CUBE_OAUTH2_ERR_TOKEN_NOT_FOUND\n"+
		"message: token not found" {
		t.Fatalf("Unexpected message:\nGot:\t\t%s\nExpected:\t%s\n", out, packet)
	}
	fmt.Println("123")
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
