# Platform independent Reverseshell and Manager in go lang
  
## Please check mask and shadow for more stealthy shells :)
### https://github.com/diljithishere/mask
### https://github.com/diljithishere/shadow

## Getting Started
### git clone https://github.com/diljithishere/reversehellgo.git
### go get github.com/fatih/color

#### cd reverseshellgo
### Build reverseshell
#### Update your manager ip and port { MANAGERIP := "ip:port" }
#### GOOS=windows GOARCH=386 go build -o hook.exe reverseshell.go (For windows executable)
#### go build -o hook reverseshell.go (Linux binary)

### Build Manager
#### go build -o manager.exe manageshell.go (Windows)
#### go build manageshell.go (Linux)

### Run manager
#### ./manageshell

##### Will get a hooked prompt after a successfull reverse connection

### Prerequisites

#### go 1.9

### Author
#### Diljith S - Initial work - (https://github.com/diljithishere)
