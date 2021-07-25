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
var now = make(map[string]time.Time)

func Check(pName string, port *serial.Port) error {
	var ret string
	log.Printf("检查设备 %s【开始】...", pName)
	if port == nil {
		_, ok := now[pName]
		if !ok || time.Now().Unix()-now[pName].Unix() > 3 {
			port, err := serial.OpenPort(&serial.Config{Name: pName, Baud: 115200, ReadTimeout: time.Second * 2})
			if err != nil {
				ret = "出错了，无法读取串口设备！" + err.Error()
				logIO.Error(ret)
				return errors.New(ret)
			} else {
				defer func() {
					err := port.Close()
					log.Println(err)
				}()
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
			now[pName] = time.Now()
		} else {
			return errors.New("您的操作太频繁了！")
		}
	}
	log.Printf("检查设备 %s【完毕】", pName)
	return nil
}
