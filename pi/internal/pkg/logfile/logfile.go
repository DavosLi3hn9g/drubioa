package logfile

import (
	"VGO/pi/internal/config"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	logPath = config.ENV["home_path"] + config.ENV["log_path"]
)

func Fatal(er interface{}) {

	log.SetPrefix("[ FATAL ]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	file, err := os.OpenFile(logPath+"error_"+time.Now().Format("200601")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}
	defer func(file *os.File) { _ = file.Close() }(file)
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	s := fmt.Sprint(er)
	_ = log.Output(2, s+" [程序被终止]")
	os.Exit(1)
}
func Error(er interface{}) {

	log.SetPrefix("[ ERROR ]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	file, err := os.OpenFile(logPath+"error_"+time.Now().Format("200601")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}
	defer func(file *os.File) { _ = file.Close() }(file)
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	s := fmt.Sprint(er)
	_ = log.Output(2, s+"")
}
func Warning(er interface{}) {

	log.SetPrefix("[ WARNING ]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	file, err := os.OpenFile(logPath+"warning_"+time.Now().Format("20060102")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}
	defer func(file *os.File) { _ = file.Close() }(file)
	log.SetOutput(io.MultiWriter(os.Stdout, file))

	s := fmt.Sprint(er)
	_ = log.Output(2, s)
}

func Info(er interface{}) {

	log.SetPrefix("[ INFO ]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	file, err := os.OpenFile(logPath+"info_"+time.Now().Format("200601")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}
	defer func(file *os.File) { _ = file.Close() }(file)
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	s := fmt.Sprint(er)
	_ = log.Output(2, s)
}
