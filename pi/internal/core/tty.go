package core

import (
	"github.com/tarm/serial"
	"io/ioutil"
	"strings"
	"time"
)

type TTY struct {
	Name string
	Desc string
}

func (tty TTY) FindAll() []TTY {
	var ttyList []TTY
	files, _ := ioutil.ReadDir("/dev/")
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "tty") {
			pName := "/dev/" + f.Name()
			ttyList = append(ttyList, TTY{Name: pName})
		}
	}
	return ttyList
}

func (tty TTY) FindUSBAudio() string {
	contents, _ := ioutil.ReadDir("/dev")
	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usbserial") ||
			strings.Contains(f.Name(), "ttyUSB") {
			return "/dev/" + f.Name()
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
