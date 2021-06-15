#!/bin/sh
go build -o ./iqiar ./pi/main.go
sudo ./iqiar $*