package sim

import (
	"github.com/tarm/serial"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type Audio struct {
	Port *serial.Port
}

func (au *Audio) Read(ch chan []byte, file chan string, usb string) {

	logIO.Print("读取数据中...")
	var pcm []byte
	var pcmPath string
	// 读缓冲区长度与data bits一致
	buf := make([]byte, 640)
	var isEnd bool
	for {
		n, err := au.Port.Read(buf)
		if err != nil {
			if err == io.EOF {
				isEnd = true
				n = 0
			} else {
				logIO.Fatal(err)
				return
			}
		}
		pcm = append(pcm, buf[:n]...)
		if n == 0 || isEnd {
			if len(pcm) == 0 {
				file <- ""
				logIO.Println("读取失败！请检查USB端口配置")
			} else {
				logIO.Print(usb + " 正在写入PCM...")
				t := time.Now().Format("20060102_15_04_05")
				filename := "call_" + t
				if configENV["pcm_path"] != "" {
					pcmPath = configENV["home_path"] + configENV["pcm_path"]
				} else {
					pcmPath = configENV["home_path"] + "data/pcm/"
				}
				err = ioutil.WriteFile(pcmPath+filename+".pcm", pcm, os.ModePerm)
				if err != nil {
					logIO.Fatal(err)
					return
				} else {
					file <- filename
					logIO.Println("...完成！")
				}
			}
		}
		bb := buf[:n]
		ch <- bb
		if n == 0 || isEnd {
			logIO.Print("读取通道关闭！")
			close(ch)
			return
		}
	}
}
func (au *Audio) Write(b []byte) {
	if _, err := au.Port.Write(b); err != nil {
		logIO.Fatal(err, "发送语音时出现严重错误！请检查语音读写端口。")
		return
	}
}
