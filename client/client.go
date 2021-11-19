package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"time"
)

const ByteSize = 4

func main() {

	flag.Parse()
	argsName := os.Args[1:] // TODO: exception len less than 4
	hostName := argsName[0]
	portName := argsName[1]

	dest := hostName + ":" + portName
	//fmt.Printf("Connecting to %s...\n", dest)

	conn, err := net.Dial("tcp", dest)

	if err != nil {
		if _, t := err.(*net.OpError); t {
			fmt.Println("Some problem connecting.")
		} else {
			fmt.Println("Unknown error: " + err.Error())
		}
		os.Exit(1)
	}

	go readConnection(conn)

	for {
		params := append([]byte(argsName[2]), byte(' '))
		params = append(params, []byte(argsName[3])...)
		params = append(params, byte('\n'))
		packet := buildPacketRequest(
			0x00000002,
			1,
			0x00000002,
			[]byte(argsName[2]),
			[]byte(argsName[3]),
		)
		packet = append(packet, byte('\n'))

		conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
		_, err = conn.Write(packet)

		if err != nil {
			fmt.Println("Error writing to stream.")
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("")
		text, _ := reader.ReadString('\n')
		text = text
	}
}

func readConnection(conn net.Conn) {
	for {
		scanner := bufio.NewScanner(conn)

		scanner.Split(bufio.ScanBytes)

		var buffer []byte
		for scanner.Scan() {
			if len(buffer) == 16 { //  header length = 12 bytes + return_code = 4 bytes
				left := 0

				var svcId int32
				var bodyLength int32
				var requestId int32
				var returnCode int32

				svcId, left = readInt32InPacket(buffer, left)
				svcId = svcId
				bodyLength, left = readInt32InPacket(buffer, left)
				requestId, left = readInt32InPacket(buffer, left)
				requestId = requestId
				returnCode, left = readInt32InPacket(buffer, left)

				//fmt.Printf("<header> ::= ")
				//fmt.Printf("%d\n", svcId)
				//fmt.Printf("%d\n", bodyLength)
				//fmt.Printf("%d\n", requestId)

				b := scanner.Bytes()
				buffer = append(buffer, b[0])

				for scanner.Scan() {
					b = scanner.Bytes()
					buffer = append(buffer, b[0])

					if int32(len(buffer)) == (bodyLength + 3*ByteSize) {
						if returnCode == 0x00000000 {
							var clientIdLen int32
							var clientId []byte
							var clientType int32
							var usernameLen int32
							var username []byte
							var expiresIn int32
							var userId int64

							clientIdLen, left = readInt32InPacket(buffer, left)
							clientId, left = readSliceBytePacket(buffer, left, clientIdLen)
							clientType, left = readInt32InPacket(buffer, left)
							usernameLen, left = readInt32InPacket(buffer, left)
							username, left = readSliceBytePacket(buffer, left, usernameLen)
							expiresIn, left = readInt32InPacket(buffer, left)
							userId, left = readInt64InPacket(buffer, left)

							//fmt.Println("<svc_ok_response_body> ::= ")
							fmt.Printf("return_сode: %d\n", returnCode)
							//fmt.Printf("<%08x>", clientIdLen)
							fmt.Printf("client_id: %s\n", clientId)
							fmt.Printf("client_type: %d\n", clientType)
							//fmt.Printf("<%08x>", usernameLen)
							fmt.Printf("username: %s\n", username)
							fmt.Printf("expires_in: %d\n", expiresIn)
							fmt.Printf("user_id: %d", userId)
							os.Exit(0)
						} else {
							var errorStringLen int32
							var errorString []byte

							errorStringLen, left = readInt32InPacket(buffer, left)
							errorString, left = readSliceBytePacket(buffer, left, errorStringLen)

							stdOutCodeError(returnCode)
							fmt.Printf("message: %s", errorString)
							os.Exit(0)
						}
					}
				}
				buffer = buffer[:0]
			}

			b := scanner.Bytes()
			buffer = append(buffer, b[0])

		}
	}
}

func handleCommands(text string) bool {
	r, err := regexp.Compile("^%.*%$")
	if err == nil {
		if r.MatchString(text) {

			switch {
			case text == "%quit%":
				fmt.Println("\b\bServer is leaving. Hanging up.")
				os.Exit(0)
			}
			return true
		}
	}
	return false
}

func readInt32InPacket(buffer []byte, left int) (int32, int) {
	integer32 := int32(binary.BigEndian.Uint32(buffer[left : left+4]))
	left = left + 4
	return integer32, left
}

func readInt64InPacket(buffer []byte, left int) (int64, int) {
	integer64 := int64(binary.BigEndian.Uint64(buffer[left : left+8]))
	left = left + 8
	return integer64, left
}

func readSliceBytePacket(buffer []byte, left int, bytesNumb int32) ([]byte, int) {
	byteArray := buffer[left : left+int(bytesNumb)]
	left = left + int(bytesNumb)
	return byteArray, left
}

func stdOutCodeError(code int32) {
	switch code {
	case 0x00000001:
		fmt.Printf("error: CUBE_OAUTH2_ERR_TOKEN_NOT_FOUND\n")
	case 0x00000002:
		fmt.Printf("error: CUBE_OAUTH2_ERR_DB_ERROR\n")
	case 0x00000003:
		fmt.Printf("error: CUBE_OAUTH2_ERR_UNKNOWN_MSG\n")
	case 0x00000004:
		fmt.Printf("error: CUBE_OAUTH2_ERR_BAD_PACKET\n")
	case 0x00000005:
		fmt.Printf("error: CUBE_OAUTH2_ERR_BAD_CLIENT\n")
	case 0x00000006:
		fmt.Printf("error: CUBE_OAUTH2_ERR_BAD_SCOPE\n")
	default:
		fmt.Printf("error: CUBE_OAUTH2_UNKNOWN_ERROR\n")
	}
}

func buildPacketRequest(
	svcId int32,
	requestId int32,

	svcMsg int32,
	token []byte,
	scope []byte,
) []byte {

	bsInt32 := make([]byte, 4)

	var packet []byte
	var body []byte

	//
	//  body filling
	//
	binary.BigEndian.PutUint32(bsInt32, uint32(svcMsg))
	body = append(body, bsInt32...)

	binary.BigEndian.PutUint32(bsInt32, uint32(len(token)))
	body = append(body, bsInt32...)
	body = append(body, token...)

	binary.BigEndian.PutUint32(bsInt32, uint32(len(scope)))
	body = append(body, bsInt32...)
	body = append(body, scope...)

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
