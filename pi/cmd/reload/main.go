package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

func main() {
	var flagPid = flag.Int("pid", 0, "进程号")
	flag.Parse()
	if *flagPid > 0 {
		cmd := exec.Command("kill", "-9", fmt.Sprintf("%d", *flagPid)) //-2 优雅终止
		if out, err := cmd.Output(); err != nil {
			log.Println(err)
			log.Println("kill failed")
		} else {
			log.Println(strings.Trim(string(out), "\n"))
			log.Println("kill success")
		}
		time.Sleep(3 * time.Second)
		cm := exec.Command("./go_build_mac", "-p", "90")
		if out, err := cm.Output(); err != nil {
			log.Println(err)
			log.Println("reload failed")
		} else {
			log.Println(strings.Trim(string(out), "\n"))
			log.Println("reload success")
		}
	} else {
		log.Println("PID cannot be empty")
	}

}
