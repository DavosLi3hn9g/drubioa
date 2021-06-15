package v1

import (
	"VGO/pi/internal/aliyun"
	"VGO/pi/internal/audio"
	"VGO/pi/internal/cons"
	"VGO/pi/internal/core"
	"VGO/pi/internal/sim"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Test struct {
}

func (t Test) FileASR(c *gin.Context) {
	filePath := c.PostForm("path")
	if filePath != "" {
		go func() {
			Ctx, Cancel := context.WithTimeout(context.Background(), 300*time.Second)
			audio.FileASR(filePath, "custom", Ctx, Cancel)
			<-Ctx.Done()
		}()
		jsonResult(c, http.StatusOK, true)
	} else {
		jsonResult(c, http.StatusBadRequest, false)
	}

}

type TestPCM struct {
	Path string `form:"path" xml:"path" json:"path"`
}

func (t Test) ListPCM(c *gin.Context) {
	var list []TestPCM
	if configENV["pcm_path"] != "" {
		files, _ := ioutil.ReadDir(configENV["home_path"] + configENV["pcm_path"])
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".pcm") && f.Size() > 100 {
				list = append(list, TestPCM{configENV["pcm_path"] + f.Name()})
			}
		}
	}
	jsonResult(c, http.StatusOK, list)
}

func (t Test) Setting(c *gin.Context) {
	var setList []Setting
	var setMap = make(map[string]string)
	var okMap = make(map[string]bool)
	type data struct {
		error  error
		status int
	}
	err := c.ShouldBindBodyWith(&setList, binding.JSON)
	if err == nil {
		for _, v := range setList {
			setMap[v.Key] = v.Value
		}
		if setMap["tty_call"] != "" {
			core.CallEnd(0, false)
			err := sim.Check(setMap["tty_call"], core.SerialAT.Port)
			if err != nil {
				jsonErr(c, http.StatusBadRequest, cons.JsonErrTTYCall, err.Error())
				return
			}
			okMap["tty_call"] = true
		}
		if setMap["tty_usb"] != "" {
			core.CallEnd(0, false)
			err := sim.Check(setMap["tty_usb"], core.SerialAudio.Port)
			if err != nil {
				jsonErr(c, http.StatusBadRequest, cons.JsonErrTTYUSB, err.Error())
				return
			}
			okMap["tty_usb"] = true
		}
		if setMap["aliyun_ak_id"] != "" && setMap["aliyun_ak_secret"] != "" {
			alConfig := &aliyun.Config{
				AccessKeyId:     setMap["aliyun_ak_id"],
				AccessKeySecret: setMap["aliyun_ak_secret"],
			}
			isi := new(aliyun.ClientISI)
			err := isi.CreateToken(alConfig)
			if err != nil {
				jsonErr(c, http.StatusBadRequest, cons.JsonErrAliyunAccess, err.Error())
				return
			}
			okMap["aliyun_ak_secret"] = true
			if setMap["aliyun_isi_appkey"] != "" {
				isi.ISIappKey = setMap["aliyun_isi_appkey"]
				if setMap["tts_voice"] != "" && setMap["tts_test"] != "" {
					isi.TestTTS(aliyun.TTSConfig{
						Voice:  setMap["tts_voice"],
						Volume: setMap["tts_volume"],
						Speech: setMap["tts_speech"],
						Pitch:  setMap["tts_pitch"],
					})

					_, err = isi.TTS(setMap["tts_test"], "test")
					if err != nil {
						jsonErr(c, http.StatusBadRequest, cons.JsonErrAliyunISIAppKey, err.Error())
						return
					}
					testWav := audio.Pcm2WavFile(configENV["home_path"]+configENV["cache_path"]+"test.pcm", "./"+configENV["wav_path"]+"test.wav")
					jsonResult(c, http.StatusOK, map[string]string{"wav": "//" + c.Request.Host + "/" + testWav.Name()})
					return

				} else {
					isi.TestTTS(aliyun.TTSConfig{
						Voice:  "xiaowei",
						Volume: "50",
						Speech: "50",
						Pitch:  "0",
					})
					_, err = isi.TTS("语音合成正常", "")
					if err != nil {
						jsonErr(c, http.StatusBadRequest, cons.JsonErrAliyunISIAppKey, err.Error())
						return
					}
				}
				okMap["aliyun_isi_appkey"] = true
			}

			if setMap["aliyun_bucket_name"] != "" {
				bucket := strings.Split(setMap["aliyun_bucket_name"], ".")
				if len(bucket) > 1 {
					bucketName := bucket[0]
					bucketLocation := bucket[1]
					oss := new(aliyun.ClientOSS)
					oss.RegionId = bucketLocation
					isExist, err := oss.NewClient(alConfig).IsBucketExist(bucketName)
					if err != nil {
						jsonErr(c, http.StatusBadRequest, cons.JsonErrAliyunOSS, err.Error())
						return
					}
					if isExist {
						okMap["aliyun_bucket_name"] = true
					} else {
						jsonErr(c, http.StatusBadRequest, cons.JsonErrAliyunOSS, "存储空间不存在！")
						return
					}
				}
			}

		}
	}
	jsonResult(c, http.StatusOK, okMap)
}
