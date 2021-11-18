/*
A very simple TCP client written in Go.

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
	"regexp"
	"time"
)

const MinInt8, MaxInt8 = -(-128), 127
const MinInt32, MaxInt32 = -(-2147483648), 2147483647
const MinInt64, MaxInt64 = -(-9223372036854775808), 9223372036854775807

const ByteSize = 4

func main() {

	flag.Parse()
	argsName := os.Args[1:] // TODO: exception len less than 4
	hostName := argsName[0]
	portName := argsName[1]

	dest := hostName + ":" + portName
	fmt.Printf("Connecting to %s...\n", dest)

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
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("")
		text, _ := reader.ReadString('\n')

		conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error writing to stream.")
			break
		}
	}
}

func readConnection(conn net.Conn) {
	for {
		scanner := bufio.NewScanner(conn)

		scanner.Split(bufio.ScanBytes)

		var buffer []byte
		for scanner.Scan() {
			if len(buffer) == 12 { //  header length = 12 bytes
				left := 0

				var svcId int32
				var bodyLength int32
				var requestId int32

				svcId, left = readInt32InPacket(buffer, left)
				svcId = svcId
				bodyLength, left = readInt32InPacket(buffer, left)
				requestId, left = readInt32InPacket(buffer, left)
				requestId = requestId

				b := scanner.Bytes()
				buffer = append(buffer, b[0])

				for scanner.Scan() {
					b = scanner.Bytes()
					buffer = append(buffer, b[0])

					if int32(len(buffer)) == (bodyLength + 3*ByteSize) {

						var returnCode int32
						var clientIdLen int32
						var clientId []byte
						var clientType int32
						var usernameLen int32
						var username []byte
						var expiresIn int32
						var userId int64

						returnCode, left = readInt32InPacket(buffer, left)
						clientIdLen, left = readInt32InPacket(buffer, left)
						clientId, left = readSliceByteInPacket(buffer, left, clientIdLen)
						clientType, left = readInt32InPacket(buffer, left)
						usernameLen, left = readInt32InPacket(buffer, left)
						username, left = readSliceByteInPacket(buffer, left, usernameLen)
						expiresIn, left = readInt32InPacket(buffer, left)
						userId, left = readInt64InPacket(buffer, left)

						//fmt.Printf("<header> ::= ")
						//fmt.Printf("%d\n", svcId)
						//fmt.Printf("%d\n", bodyLength)
						//fmt.Printf("%d\n", requestId)

						//fmt.Println("<svc_ok_response_body> ::= ")
						fmt.Printf("return_сode: %d\n", returnCode)
						fmt.Printf("<%08x>", clientIdLen)
						fmt.Printf("client_id: %s\n", clientId)
						fmt.Printf("client_type: %d\n", clientType)
						//fmt.Printf("<%08x>", usernameLen)
						fmt.Printf("username: %s\n", username)
						fmt.Printf("expires_in: %d\n", expiresIn)
						fmt.Printf("user_id: %d\n", userId)
					}
				}
				buffer = buffer[:0]
			}

			// Get Bytes and display the byte.
			b := scanner.Bytes()
			buffer = append(buffer, b[0])

		}

		for {
			ok := scanner.Scan()
			text := scanner.Text()

			command := handleCommands(text)
			if !command {
				fmt.Printf("%s\n> ", text)
			}

			if !ok {
				fmt.Println("Reached EOF on server connection.")
				break
			}
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

func readSliceByteInPacket(buffer []byte, left int, bytesNumb int32) ([]byte, int) {
	byteArray := buffer[left : left+int(bytesNumb)]
	left = left + int(bytesNumb)
	return byteArray, left
}
