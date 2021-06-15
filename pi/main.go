package main

import (
	"VGO/pi/internal/cache"
	"VGO/pi/internal/cmd"
	"VGO/pi/internal/config"
	"VGO/pi/internal/core"
	"VGO/pi/internal/orm"
	"VGO/pi/internal/pkg/file"
	"VGO/pi/internal/routers"
	v1 "VGO/pi/internal/routers/api/v1"
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	orm.OpenDB()

}
func main() {

	var configENV = config.ENV
	var err error
	var host, port string
	var logIO = new(file.Log)
	var (
		flagPW = flag.Bool("reset", false, "重置控制台密码")
		flagP  = flag.String("p", "", "API服务端口号")
		flagD  = flag.Bool("d", false, "是否开启Debug")
		flagR  = flag.String("r", "custom", "turing：图灵机器人回复 \r custom：自定义回复 ")
		flagT  = flag.String("t", "", "用于测试的本地文件路径  data/pcm/tonghua_20190804_10_46.pcm")
	)
	flag.Parse()

	core.FlagReply = flagR
	core.FlagTest = flagT

	if *flagPW {
		vAuth := new(v1.Auth)
		newPW := vAuth.Reset()
		fmt.Println("==========")
		fmt.Println("密码已经重置为：", newPW)
		fmt.Println("请进入控制台>系统设置修改密码。")
		fmt.Println("==========")
		pause()
	}
	core.Setting = cache.SettingsCache.Read()
	core.PowerOn(5 * time.Second)
	go core.Start()
	go cmd.LoopNetwork()
	var router *gin.Engine

	if *flagP != "" {
		port = *flagP
	} else {
		port = configENV["port"]
	}
	if *flagD || configENV["dev_mode"] == "true" {
		router = routers.List(true)
	} else {
		router = routers.List(false)
	}
	host = ":" + port
	srv := &http.Server{Addr: host, Handler: router}
	cmd.Pid = syscall.Getpid()
	logIO.Printf("pid is %d\r\n", syscall.Getpid())
	go func() {
		logIO.Printf("host: %s", host)
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logIO.Printf("server error: %s\r\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logIO.Println("Shutdown Server ...")
	*core.IsRunning = false
	core.CallEnd(0, true)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logIO.Fatal("Server Shutdown:", err)
	}
	logIO.Printf("Server on %s stopped", host)
}
func pause() {
	fmt.Print("请输入回车继续...")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "" {
			break
		}
	}
}
