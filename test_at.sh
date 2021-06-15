#!/bin/sh
go build -o ./test_at ./pi/cmd/serial_at/main.go
sudo ./test_at $*