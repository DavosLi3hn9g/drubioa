package cmd

import (
	"VGO/pi/internal/pkg/file"
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

var basePath = "/bin/bash"
var logIO *file.Log
var Pid = 0

func Bash(sh string) string {
	var err error
	var cmd *exec.Cmd
	var out []byte

	cmd = exec.Command(basePath, "-c", sh)
	//阻塞执行
	if out, err = cmd.Output(); err != nil {
		logIO.Error(err)
	}
	return strings.Trim(string(out), "\n")
}

func Shell(sh string, timeout time.Duration) {
	var err error
	var cmd *exec.Cmd
	var line string
	if timeout > 0 {
		//超时终止进程
		ctx, cancel := context.WithTimeout(context.Background(), (timeout)*time.Second)
		defer cancel()
		cmd = exec.CommandContext(ctx, basePath, "-c", sh)
	} else {
		cmd = exec.Command(basePath, "-c", sh)
	}

	fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	//非阻塞执行
	_ = cmd.Start()
	//一行一行读取
	reader := bufio.NewReader(stdout)
	for {
		line, err = reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		fmt.Println(line)
	}
	//阻塞，等待命令返回后释放相关的资源
	_ = cmd.Wait()
	defer stdout.Close()
}
