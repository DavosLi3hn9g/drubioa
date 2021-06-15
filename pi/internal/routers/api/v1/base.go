package v1

import (
	"VGO/pi/internal/config"
	"VGO/pi/internal/pkg/file"
	"github.com/gin-gonic/gin"
)

var logIO *file.Log
var configENV = config.ENV

func jsonResult(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{
		"result": data,
	})
}

func jsonErr(c *gin.Context, code int, err int, msg string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": msg,
		"error":   err,
	})
}
