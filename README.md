# adb go

Executable Go package to make connecting/pairing Android devices through adb more interactive.

## Installation

```shell
go install github.com/dnieln7/adb-go@latest
``` 

## Usage

1. Enable wireless debugging and take note of your ip and port.
2. Run `adb-go`.
3. Select your device's ip:port.

* If the connection fails the program will ask for a port and code, go to "Pair device with pairing code" in the
  wireless debugging screen. After pairing, you can try to connect again. 
