package core

import (
	"github.com/tarm/serial"
	"go.bug.st/serial/enumerator"
	"io/ioutil"
	"strings"
	"time"
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
	files, _ := ioutil.ReadDir("/dev/")
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "tty") {
			// it is a legitimate serial port
			pname := "/dev/" + f.Name()
			c2 := &serial.Config{Name: pname, Baud: 115200, ReadTimeout: time.Second * 5}
			_, err := serial.OpenPort(c2)
			if err != nil {
				//core.Println(err.Error())
			} else {
				logIO.Println(pname)
			}

		}

	}
}
