package sim

import (
	"VGO/pi/internal/config"
	"VGO/pi/internal/pkg/file"
	"errors"
	"github.com/tarm/serial"
	"io"
	"log"
	"time"
)

var logIO *file.Log
var configENV = config.ENV

func Check(pName string, port *serial.Port) error {
	var ret string
	log.Printf("检查设备 %s【开始】...", pName)
	if port == nil {
		port, err := serial.OpenPort(&serial.Config{Name: pName, Baud: 115200, ReadTimeout: time.Second * 2})
		if err != nil {
			ret = "出错了，无法读取串口设备！" + err.Error()
			logIO.Error(ret)
			return errors.New(ret)
		} else {
			defer port.Close()
			buf := make([]byte, 640)
			n, err := port.Read(buf)
			if err != nil {
				if err == io.EOF {
					n = 0
				} else {
					panic(err)
				}
			}
			log.Printf("读取：%d %v", n, buf[:n])
		}
	}
	log.Printf("检查设备 %s【完毕】", pName)
	return nil
}
