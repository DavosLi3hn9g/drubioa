package core

import (
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type TTY struct {
	Name  string
	Desc  string
	Error string
	IsUSB bool
}

func (tty TTY) FindAll() []TTY {
	var ttyList []TTY
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		logIO.Fatal(err)
	}
	if len(ports) == 0 {
		logIO.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		ttyList = append(ttyList, TTY{Name: port.Name, IsUSB: port.IsUSB, Desc: port.Product})
	}
	return ttyList
}

func (tty TTY) FindUSBAudio() string {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		logIO.Fatal(err)
	}
	if len(ports) == 0 {
		logIO.Println("No serial ports found!")
		return ""
	}
	for _, port := range ports {
		logIO.Printf("Found port: %s\n", port.Name)
		if port.IsUSB {
			logIO.Printf("USB ID %s:%s\n", port.VID, port.PID)
			logIO.Printf("USB serial %s\n", port.SerialNumber)
		}
	}
	return ""
}

func (tty TTY) DevLoop() {
	for _, f := range tty.FindAll() {
		// it is a legitimate serial port
		pname := f.Name
		_, err := serial.Open(pname, &serial.Mode{BaudRate: 115200})
		if err != nil {
			//core.Println(err.Error())
		} else {
			logIO.Println(pname)
		}

	}
}
