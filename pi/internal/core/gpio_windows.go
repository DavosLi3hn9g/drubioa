package core

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
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
	return
}
func LoopPin(BCM uint8, ch chan string) {
	return
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
