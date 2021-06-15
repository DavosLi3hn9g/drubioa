package v1

import (
	"VGO/pi/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/tarm/serial"
	"net/http"
	"strings"
	"time"
)

type Device struct{}

func (dv Device) ListTTY(c *gin.Context) {
	var listSerial []core.TTY
	var listUSB []core.TTY
	var list = make(map[string]interface{}, 2)
	listTTY := core.TTY{}.FindAll()
	for _, v := range listTTY {
		if strings.HasPrefix(v.Name, "/dev/ttyAMA") || strings.HasPrefix(v.Name, "/dev/ttyS") {
			c2 := &serial.Config{Name: v.Name, Baud: 115200, ReadTimeout: time.Second * 5}
			p, err := serial.OpenPort(c2)
			if err != nil {
				v.Desc = "不可用"
			} else {
				v.Desc = ""
				_ = p.Flush()
				_ = p.Close()
			}

			listSerial = append(listSerial, v)
		}
		if strings.HasPrefix(v.Name, "/dev/ttyUSB") {
			listUSB = append(listUSB, v)
		}

	}
	list["serial"] = listSerial
	list["usb"] = listUSB
	jsonResult(c, http.StatusOK, list)
}
