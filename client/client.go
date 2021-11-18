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
				var left *int
				var right *int
				tmpLeft, tmpRight := 0, 4

				left, right = &tmpLeft, &tmpRight

				//svcId := int32(binary.BigEndian.Uint32(buffer[left:right]))
				*left, *right = *left+ByteSize, *right+ByteSize

				bodyLength := int32(binary.BigEndian.Uint32(buffer[*left:*right]))
				*left, *right = *left+ByteSize, *right+ByteSize

				//requestId := int32(binary.BigEndian.Uint32(buffer[left:right]))
				*left, *right = *left+ByteSize, *right+ByteSize

				b := scanner.Bytes()
				buffer = append(buffer, b[0])

				for scanner.Scan() {
					b = scanner.Bytes()
					buffer = append(buffer, b[0])

					if int32(len(buffer)) == (bodyLength + 3*ByteSize) {

						var returnCode *int32
						//int32(binary.BigEndian.Uint32(buffer[left:right]))
						readInt32InPacket(returnCode, buffer, left, right)
						//left, right = left+ByteSize, right+ByteSize
						//
						//clientIdLen := int32(binary.BigEndian.Uint32(buffer[left:right]))
						//left, right = left+ByteSize, right+int(clientIdLen)
						//
						//clientId := buffer[left:right]
						//left, right = left+int(clientIdLen), right+ByteSize
						//
						//clientType := int32(binary.BigEndian.Uint32(buffer[left:right]))
						//left, right = left+ByteSize, right+ByteSize
						//
						//usernameLen := int32(binary.BigEndian.Uint32(buffer[left:right]))
						//left, right = left+ByteSize, right+int(usernameLen)
						//
						//username := buffer[left:right]
						//left, right = left+int(usernameLen), right+ByteSize
						//
						//expiresIn := int32(binary.BigEndian.Uint32(buffer[left:right]))
						//left, right = left+ByteSize, right+(2*ByteSize)
						//
						//userId := int64(binary.BigEndian.Uint64(buffer[left:right]))

						//fmt.Printf("<header> ::= ")
						//fmt.Printf("%d\n", svcId)
						//fmt.Printf("%d\n", bodyLength)
						//fmt.Printf("%d\n", requestId)

						//fmt.Println("<svc_ok_response_body> ::= ")
						fmt.Printf("return_сode: %d\n", returnCode)
						//fmt.Printf("<%08x>", clientIdLen)
						//fmt.Printf("client_id: %s\n", clientId)
						//fmt.Printf("client_type: %d\n", clientType)
						////fmt.Printf("<%08x>", usernameLen)
						//fmt.Printf("username: %s\n", username)
						//fmt.Printf("expires_in: %d\n", expiresIn)
						//fmt.Printf("user_id: %d\n", userId)
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

func readInt32InPacket(int32Var *int32, buffer []byte, left *int, right *int) {
	integer := int32(binary.BigEndian.Uint32(buffer[*left:*right]))
	int32Var = &integer
	*left, *right = *left+4, *right+4
}
