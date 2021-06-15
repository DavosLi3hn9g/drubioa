package cmd

import (
	"time"
)

func CheckNetwork() {
	var sh = `/sbin/ifconfig | grep 192 |  awk '{print $2}'`
	if Bash(sh) == "" {
		logIO.Warning("WIFI掉线了，已重连！")
		Bash(`sudo /etc/init.d/networking restart`)
	}

}

func LoopNetwork() {
	for {
		CheckNetwork()
		time.Sleep(time.Minute * 5)
	}
}
