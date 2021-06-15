package v1

import (
	"VGO/pi/internal/aliyun"
	"VGO/pi/internal/cons"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type Cloud struct {
}

func (cl Cloud) ListOSS(c *gin.Context) {
	var setList []Setting
	var setMap = make(map[string]string)
	err := c.ShouldBindBodyWith(&setList, binding.JSON)

	if err == nil {
		for _, v := range setList {
			setMap[v.Key] = v.Value
		}
		if setMap["aliyun_ak_id"] != "" && setMap["aliyun_ak_secret"] != "" {
			oss := new(aliyun.ClientOSS)
			listBuckets, err := oss.NewClient(&aliyun.Config{
				AccessKeyId:     setMap["aliyun_ak_id"],
				AccessKeySecret: setMap["aliyun_ak_secret"],
			}).GetListBuckets()
			if err != nil {
				jsonErr(c, http.StatusBadRequest, cons.JsonErrAliyunOSS, err.Error())
				return
			}
			jsonResult(c, http.StatusOK, listBuckets)
			return
		}
	} else {
		jsonErr(c, http.StatusBadRequest, cons.JsonErrAliyunOSS, err.Error())
		return
	}
}
func (cl Cloud) ListTTS(c *gin.Context) {
	list := aliyun.VoiceList
	jsonResult(c, http.StatusOK, list)
}
