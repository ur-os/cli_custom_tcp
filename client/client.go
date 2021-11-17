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
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"time"
)

//var host = flag.String("host", "localhost", "The hostname or IP to connect to; defaults to \"localhost\".")
//var port = flag.Int("port", 8000, "The port to connect to; defaults to 8000.")

type IpProtoHeader struct {
	svcId int32
	bodyLength int32
	requestId int32
}

type IpProtoBody struct {
	strLength int32
	str string  // TODO: check is this correct interpretation of "int8+" from technical task
}

type IpProtoPacket struct {
	header IpProtoHeader
	body IpProtoBody
}


func main() {
	flag.Parse()
	argsName := os.Args[1:]  // TODO: exception len less than 4
	hostName := argsName[0]
	portName := argsName[1]


	dest := hostName + ":" + portName
	fmt.Printf("Connecting to %s...\n", dest)

	conn, err := net.Dial("tcp", dest)

	if err != nil {
		if _, t := err.(*net.OpError); t {
			fmt.Println("Сonnection not established")
		} else {
			fmt.Println("Unknown error: " + err.Error())
		}
		os.Exit(1)
	}

	go readConnection(conn)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
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
		buf := new(bytes.Buffer)
		binary.Read(conn, binary.LittleEndian, buf)

		for {
			ok := scanner.Scan()
			text := scanner.Text()
			bytes := scanner.Bytes()

			command := handleCommands(text, bytes)
			if !command {
				fmt.Printf("%s\n", text)
			}

			if !ok {
				fmt.Println("Reached EOF on server connection.")
				break
			}
		}
	}
}

func handleCommands(text string, bytes []byte) bool {
	r, err := regexp.Compile("^%.*%$")
	result := regexp.MustCompile(`\n`)
	//resultBytes := bytes



	fmt.Println(result.Split(text, -1))
	if err == nil {
		if r.MatchString(text) {

			switch {
			case text == "%quit%":
				fmt.Println("\b\bServer is leaving. Hanging up.")
				os.Exit(0)
			}

			fmt.Println("handleCommands end")
			return true
		}
	}
	fmt.Println("handleCommands end")
	return false
}

