package core

import (
	"VGO/pi/internal/cmd"
	"encoding/json"
	"github.com/stianeikeland/go-rpio/v4"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var PinJson GPin
var PinCache = make([]uint8, 0, 1)

type GPin struct {
	Outside []struct {
		BCM string
	}
	Inside []struct {
		BCM string
	}
}

func PowerOn(wait time.Duration) {
	var sh = `uname -a | grep raspberrypi |  awk '{print $2}'`
	if cmd.Bash(sh) != "" {
		if err := rpio.Open(); err != nil {
			logIO.Fatal(err)
		} else {
			defer rpio.Close()
		}
		var pin = make(map[int]rpio.Pin, 2)
		powerPin := strings.Split(configENV["bcm_power"], ",")
		if len(powerPin) == 2 {
			for k, v := range powerPin {
				p, err := strconv.Atoi(v)
				if err != nil {
					logIO.Fatal("bcm_power配置异常！", err)
				} else if p > 0 {
					pin[k] = rpio.Pin(p)
					pin[k].Output()
					pin[k].Low()
					time.Sleep(100 * time.Millisecond)
				}
			}
			logIO.Println("正在启动SIM设备")
			time.Sleep(wait)
		} else {
			logIO.Fatal("bcm_power未配置，SIM设备无法自动开机！")
		}
	}
}
func LoopPin(BCM uint8, ch chan string) {
	if err := rpio.Open(); err != nil {
		logIO.Fatal(err)
	} else {
		defer rpio.Close()
		var pin = map[uint8]rpio.Pin{}
		var state = map[uint8]rpio.State{}
		var old = map[uint8]rpio.State{}

		for {
			if len(PinCache) == 0 {
				pin[BCM] = rpio.Pin(BCM)
				state[BCM] = pin[BCM].Read()
				PinCache = append(PinCache, BCM)
			}
			for _, i := range PinCache {
				old[i] = state[i]
				state[i] = pin[i].Read()
				if state[i] != old[i] {
					if i == BCM {
						if state[i] == 0 {
							ch <- "START"
							logIO.Println("BCM LOW")
						} else {
							//ch <- "END"
							logIO.Println("BCM HIGH")
						}
					}
					logIO.PinState(i, uint8(state[i]))
				}
			}
			time.Sleep(time.Second * 2) //读取频率
		}
	}
}
func PinCacheRead() GPin {
	jsonFile := "./configs/gpio.json"
	p, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		logIO.Error("找不到GPIO配置文件: %v", err)
	}
	err = json.Unmarshal(p, &PinJson)
	if err != nil {
		logIO.Error("GPIO配置错误: %v", err)
	}
	PinCacheClear()
	var PinAll []struct {
		BCM string
	}
	PinAll = append(PinJson.Inside, PinJson.Outside...)
	for _, v := range PinAll {
		if v.BCM != "*" {
			iv, err := strconv.ParseUint(v.BCM, 0, 8)
			if err == nil {
				PinCache = append(PinCache, uint8(iv))
			}
		}
	}
	return PinJson
}
func PinCacheClear() {
	PinCache = nil
}
