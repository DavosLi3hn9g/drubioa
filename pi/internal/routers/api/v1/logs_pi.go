package v1

import (
	"VGO/pi/internal/core"
	"VGO/pi/internal/pkg/file"
	"github.com/gin-gonic/gin"
	"github.com/tarm/serial"
	"net/http"
)

type LogsPi struct {
}

func (pi *LogsPi) port() {
	var err error
	if core.SerialAT.Port == nil {
		core.SerialAT.Port, err = serial.OpenPort(&serial.Config{Name: core.Setting.TTYCall, Baud: 115200})
		if err != nil {
			logIO.Fatal(err)
			return
		}
	}
}
func (pi *LogsPi) GPIO(c *gin.Context) {
	data := core.PinCacheRead()
	jsonResult(c, http.StatusOK, data)
}

func (pi *LogsPi) SysCache(c *gin.Context) {
	reader := file.Cache
	if c.Query("ac") == "read" {
		*file.Cache = nil
	}
	jsonResult(c, http.StatusOK, map[string]interface{}{
		"list":   reader,
		"unread": len(*reader),
	})
}

func (pi *LogsPi) All(c *gin.Context) {
	pi.port()
	data := string(logIO.ReadFile())
	jsonResult(c, http.StatusOK, data)
}
func (pi *LogsPi) Add(c *gin.Context) {
	at := c.PostForm("at")
	pi.port()
	core.SerialAT.AT(at)
	jsonResult(c, http.StatusOK, true)
}

func (pi *LogsPi) Clear(c *gin.Context) {
	logIO.Clear()
	jsonResult(c, http.StatusOK, true)
}
