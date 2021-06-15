package main

import (
	"VGO/pi/internal/core"
	"VGO/pkg/fun"
	"bytes"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"time"
)

func main() {

	var (
		b      []byte
		chGPIO = make(chan string)
		err    error
	)
	core.SerialAT.Port, err = serial.OpenPort(&serial.Config{Name: "/dev/ttyAMA0", Baud: 115200, Size: 8})
	//Port.Flush()
	defer core.SerialAT.Port.Close()
	log.Println("start serial")
	if err != nil {
		log.Fatal("serial.OpenPort ", err)
		return
	}
	go core.LoopPin(27, chGPIO)
	go StatusRead(core.SerialAT.Port) //查询状态
	//go core.PhoneCore.PiInitAT() //初始化

	for {
		var buf = make([]byte, 8)
		if core.SerialAT.Port == nil {
			log.Println("WatchAT，没有连接，正在重连...")
			time.Sleep(5 * time.Second)
			continue
		}
		n, err := core.SerialAT.Port.Read(buf) //没有任何AT数据读取时会阻塞在这
		if err != nil {

			log.Println("WatchAT:", err)
			return
		}
		//log.Printf("buf: n %d, len %d, cap %d, %v %q", n, len(buf), cap(buf), buf[:n], buf[:n])
		b = append(b, buf[:n]...)
		if bytes.HasSuffix(b, []byte{13, 10}) || bytes.HasSuffix(b, []byte{10, 10}) { //换行结尾 [13 10],注意读取SMS时的换行不一定是结尾
			data := string(b)
			data = fun.PregReplace("\r\r\n", "\r\n", data)
			b = b[:0]
			fmt.Println("----- AT -----\r\n", data)

		}

	}

}

func AT(port *serial.Port, cmd string) {
	if _, err := port.Write([]byte(cmd + "\r\n")); err != nil {
		log.Fatal("AT:", err)
	}
}

func StatusRead(Port *serial.Port) {
	wait := 3 * time.Second
	//AT(Port, "ATE1")      //关闭回显
	AT(Port, "AT+CLIP?")  //来电显示
	time.Sleep(wait)      //这里必须要停顿一下，否则会出错
	AT(Port, "AT+CPOWD?") //GSM状态
	time.Sleep(wait)
	AT(Port, "AT+CPCMFRM?") //16k
	time.Sleep(wait)
	AT(Port, "AT+CLVL?") //最大音量
	time.Sleep(wait)
	//接收中文短信
	AT(Port, "AT+CMGF?") //设置文本显示
	time.Sleep(wait)
	AT(Port, `AT+CSCS?`) //设置GSM编码集
	time.Sleep(wait)
	AT(Port, "AT+CSMP?") //设置文本模式参数
	time.Sleep(wait)
	AT(Port, "AT+CNMI?") //设置新信息提醒
	time.Sleep(wait)
	//AT(Port, `AT+CMGL="ALL"`)
}
