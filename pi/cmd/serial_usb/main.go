package main

import (
	"VGO/pi/internal/core"
	"VGO/pi/internal/sim"
	"flag"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"strings"
	"time"
)

func main() {
	var (
		flagAT    = flag.String("at", "/dev/ttyAMA0", "AT设备端口号")
		flagAudio = flag.String("au", "/dev/ttyUSB4", "Audio设备端口号，可以用英文','号间隔")
		flagTel   = flag.String("t", "10010", "测试电话")
		err       error
		ret       string
	)
	core.PowerOn(5 * time.Second)
	flag.Parse()
	SerialAT := new(sim.AT)
	SerialAT.Port, err = serial.OpenPort(&serial.Config{Name: *flagAT, Baud: 115200, ReadTimeout: time.Second * 2})
	if err != nil {
		ret = "出错了，无法读取AT串口设备！" + err.Error()
		log.Println(ret)
	}
	log.Printf("写入：%s", *flagAT)
	SerialAT.AT(fmt.Sprintf("ATD%s;", *flagTel))
	log.Printf("正在拨打：%s", *flagTel)
	time.Sleep(time.Second * 5)
	SerialAT.AT("AT+CPCMREG=1")

	postArr := strings.Split(*flagAudio, ",")

	for _, v := range postArr {
		port, err := serial.OpenPort(&serial.Config{Name: v, Baud: 115200, ReadTimeout: time.Second * 2})
		if err != nil {
			ret = "出错了，无法读取Audio串口设备！" + err.Error()
			log.Println(ret)
		} else {
			buf := make([]byte, 640)
			n, err := port.Read(buf)
			if err != nil {
				if err == io.EOF {
					n = 0
				} else {
					panic(err)
				}
			}
			log.Printf("读取：%s %d %s", v, n, buf[:n])
			port.Close()
		}
		time.Sleep(time.Second * 2)
	}
	time.Sleep(time.Second * 3)
	SerialAT.AT("AT+CPCMREG=0")
	time.Sleep(time.Second * 3)
	SerialAT.AT("AT+CHUP")
	log.Println("已挂机，测试结束！")
}
