#!/bin/sh
go build -o ./test_usb ./pi/cmd/serial_usb/main.go
sudo ./test_usb $*