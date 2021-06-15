package file

import (
	"VGO/pi/internal/config"
	"VGO/pi/internal/pkg/logfile"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Log struct {
	WsOut *io.PipeWriter
}

type Msg struct {
	T int
	M interface{}
	D string
}

var (
	buf       = make([]byte, 0, 64)
	logIO     = new(Log)
	configENV = config.ENV
	Cache     = new([]Msg)
	logPath   = configENV["home_path"] + configENV["log_path"]
)

func openFile() (*os.File, error) {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	file, err := os.OpenFile(logPath+"pi.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}
	return file, err
}
func saveFile(s []byte) {
	file, _ := openFile()
	defer func(file *os.File) { _ = file.Close() }(file)
	_, _ = file.Write(s)
	//log.Print(s)
}

func saveBuffer(s []byte, writerFile bool) {
	capSize := cap(buf)
	buf = append(buf, s...)
	if len(buf) >= capSize || writerFile {
		saveFile(buf)
		buf = make([]byte, 0, capSize)
	}
	fmt.Print(string(s))
}

func WsWrite(writer *io.PipeWriter) {
	logIO.WsOut = writer
}
func WsClose() {
	logIO.WsOut = nil
}
func (l *Log) Fatal(str ...interface{}) {
	var strP string
	for _, arg := range str {
		strP = strP + fmt.Sprint(arg)
	}
	b := "[ 严重错误 ]" + strP + "\r\n"
	if logIO.WsOut != nil {
		push, _ := json.Marshal(&Msg{
			T: 1,
			M: b,
		})
		_, err := logIO.WsOut.Write(push)
		if err != nil {
			logfile.Error(err)
		}
	}
	*Cache = append(*Cache, Msg{
		T: 1,
		M: b,
		D: time.Now().Format("2006-01-02 15:04:05"),
	})
	saveBuffer([]byte(b), false)
}
func (l *Log) Warning(str ...interface{}) {
	l.Println("[ WARNING ]", str)
}
func (l *Log) Error(str ...interface{}) {
	l.Println("[ ERROR ]", str)
}
func (l *Log) Println(str ...interface{}) {
	var strP string
	for _, arg := range str {
		strP = strP + fmt.Sprint(arg)
	}
	b := strP + "\r\n"
	if logIO.WsOut != nil {
		push, _ := json.Marshal(&Msg{
			T: 1,
			M: b,
		})
		//log.Println("logIO.WsOut")
		_, err := logIO.WsOut.Write(push)
		if err != nil {
			logfile.Error(err)
		}
	}
	saveBuffer([]byte(b), false)
}
func (l *Log) Printf(format string, a ...interface{}) {
	b := fmt.Sprintf(format, a)
	if logIO.WsOut != nil {
		push, _ := json.Marshal(&Msg{
			T: 1,
			M: b,
		})
		_, err := logIO.WsOut.Write(push)
		if err != nil {
			logfile.Error(err)
		}
	}
	saveBuffer([]byte(b), false)

}
func (l *Log) PrintlnDev(str ...interface{}) {
	if configENV["dev_mode"] == "true" {
		l.Println(str...)
	}
}
func (l *Log) PrintDev(a string) {
	if configENV["dev_mode"] == "true" {
		l.Print(a)
	}
}
func (l *Log) Print(a string) {
	b := fmt.Sprint(a)

	if logIO.WsOut != nil {
		push, _ := json.Marshal(&Msg{
			T: 1,
			M: b,
		})
		_, err := logIO.WsOut.Write(push)
		if err != nil {
			logfile.Error(err)
		}
	}
	saveBuffer([]byte(b), false)
}
func (l *Log) Write(a string) {
	b := fmt.Sprintln(a)
	if logIO.WsOut != nil {
		push, _ := json.Marshal(&Msg{
			T: 1,
			M: b,
		})
		_, err := logIO.WsOut.Write(push)
		if err != nil {
			logfile.Error(err)
		}
	}
	saveBuffer([]byte(b), true)
}

func (l *Log) Clear() {
	err := os.Truncate(logPath+"pi.log", 0)
	if err != nil {
		log.Fatalln("删除日志文件失败：", err)
	}
}
func (l *Log) ReadFile() []byte {
	filename := logPath + "pi.log"
	fileInfo, _ := os.Stat(filename)
	fileSize := fileInfo.Size()
	if fileSize > 1024*1024 {
		return []byte("日志文件太大了，请手动清理一次！")
	}
	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	return p
}
