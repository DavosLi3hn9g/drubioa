package v1

import (
	"VGO/pi/internal/core"
	"github.com/gin-gonic/gin"
	"go.bug.st/serial"
	"net/http"
	"strings"
)

type Device struct{}

func (dv Device) ListTTY(c *gin.Context) {
	var listSerial []core.TTY
	var listUSB []core.TTY
	var list = make(map[string]interface{}, 2)
	listTTY := core.TTY{}.FindAll()
	for _, v := range listTTY {
		if strings.HasPrefix(v.Name, "/dev/ttyAMA") || strings.HasPrefix(v.Name, "/dev/ttyS") || strings.HasPrefix(v.Name, "COM") {
			p, err := serial.Open(v.Name, &serial.Mode{BaudRate: 115200})
			if err != nil {
				v.Error = err.Error()
			} else {
				_ = p.Close()
			}
			listSerial = append(listSerial, v)
		}
		if strings.HasPrefix(v.Name, "/dev/ttyUSB") || v.IsUSB {
			listUSB = append(listUSB, v)
		}

	}
	list["serial"] = listSerial
	list["usb"] = listUSB
	jsonResult(c, http.StatusOK, list)
}
