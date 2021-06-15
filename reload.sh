#!/bin/sh
go build -o ./reload ./pi/cmd/reload/main.go
sudo ./reload $*