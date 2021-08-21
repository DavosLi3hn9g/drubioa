package v1

import (
	"VGO/pi/internal/cmd"
	"VGO/pi/internal/cons"
	"VGO/pi/internal/pkg/file"
	"VGO/pi/internal/pkg/logfile"
	"VGO/pkg/curl"
	"VGO/pkg/fun"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/go-update"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Version struct {
	Vid      int    `form:"vid" json:"vid" xml:"vid" gorm:"primary_key"` //
	Platform string `form:"platform" json:"platform" xml:"platform"`     //平台
	Version  string `form:"version" json:"version" xml:"version"`        //版本
	Content  string `form:"content" json:"content" xml:"content"`        //正文
	Download string `form:"download" json:"download" xml:"download"`     //下载地址
	Status   int    `form:"status" json:"status" xml:"status"`           //状态 -1开发版，0正式版，1公测版
	Dateline int    `form:"dateline" json:"dateline" xml:"dateline"`     //发布时间
	DownNum  int    `form:"down_num" json:"down_num" xml:"down_num"`     //下载量
	DateStr  string `form:"date_str" json:"date_str" xml:"date_str"`
	Hash     string `form:"hash" json:"hash" xml:"hash"`
}

var progress float64
var progressPrev float64
var progressRepeat int
var lastHash string

func (ver Version) versionLast() Version {
	var lastResp struct {
		Result Version `json:"result"`
	}
	last := curl.GET(fmt.Sprintf("https://api.iqiar.com/api/v1/ai/version_last?goos=%s&goarch=%s", runtime.GOOS, runtime.GOARCH), map[string]string{})
	_ = json.Unmarshal(last, &lastResp)
	return lastResp.Result
}

func (ver Version) GetVersion(c *gin.Context) {

	var data = make(map[string]interface{})
	data["version"] = cons.Version
	data["content"] = "当前已经是最新版本"
	data["dateline"] = fun.StrTimestamp(cons.VerDate)
	data["date_str"] = cons.VerDate
	data["last"] = nil
	lastResp := ver.versionLast()
	if lastResp.Dateline > data["dateline"].(int) {
		data["last"] = lastResp
		lastHash = lastResp.Hash
	}

	jsonResult(c, http.StatusOK, data)
}

func (ver Version) DoDownload(c *gin.Context) error {
	lastResp := ver.versionLast()
	resp, err := http.Get(lastResp.Download)
	if err != nil {
		return err
	} else {
		logIO.Println("获取下载地址")
	}
	defer resp.Body.Close()
	temp := "./update.zip"
	reader := &Reader{
		Reader: resp.Body,
		Total:  resp.ContentLength,
	}
	fp, err := os.OpenFile(temp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = io.Copy(fp, reader)
	if err != nil {
		return err
	} else {
		if lastHash != "" {
			h := sha1.New()
			_, err = io.Copy(h, reader)
			if err != nil {
				return err
			} else if fmt.Sprintf("%x", h.Sum(nil)) == lastHash {
				if err = file.Unzip(temp, ".", []string{configENV["db_path"]}); err != nil {
					return err
				}
				if err = os.Remove(temp); err != nil {
					return err
				}
				logIO.Println("系统更新完毕！")
			} else {
				return errors.New("文件验证失败，请联系开发者。")
			}
		} else {
			return errors.New("Hash异常，请检查API是否正确！")
		}
	}
	return nil
}

type Reader struct {
	io.Reader
	Total   int64
	Current int64
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	r.Current += int64(n)
	progressPrev := progress
	progress = float64(r.Current*10000/r.Total) / 100
	if progressPrev == progress {

	}
	fmt.Printf("\r进度 %.2f%%", progress)

	return
}

func (ver Version) DoApply(c *gin.Context) error {
	lastResp := ver.versionLast()
	resp, err := http.Get(lastResp.Download)
	if err != nil {
		logfile.Error("下载失败！" + err.Error())
		return err
	}
	defer resp.Body.Close()
	reader := &Reader{
		Reader: resp.Body,
		Total:  resp.ContentLength,
	}
	err = update.Apply(reader, update.Options{TargetMode: 0777})
	if err != nil {
		logfile.Error("更新失败！" + err.Error())
		return err
	}
	return nil
}

func (ver Version) Progress(c *gin.Context) {
	var data = make(map[string]interface{})
	data["progress"] = progress
	data["completed"] = false
	if progress >= 100 {
		data["completed"] = true
	} else if progress > 0 {
		if progressPrev != progress {
			progressPrev = progress
		} else {
			progressRepeat++
			if progressRepeat >= 30 {
				progress = 0
				progressPrev = 0
				jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "更新中断！请重试。")
				return
			}
		}
	}
	jsonResult(c, http.StatusOK, data)
	return
}
func (ver Version) DoUpdate(c *gin.Context) {
	go func() {
		err := ver.DoDownload(c)
		if err != nil {
			logIO.Error(err.Error())
			jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, err.Error())
		}
	}()
	var data = make(map[string]interface{})
	data["version"] = cons.Version
	data["content"] = "开始更新，请不要断网。"
	data["dateline"] = fun.StrTimestamp(cons.VerDate)
	data["date_str"] = cons.VerDate
	data["last"] = nil
	logfile.Info(data["content"])
	jsonResult(c, http.StatusOK, data)
}

func (ver Version) Reload(c *gin.Context) {

	go func() {
		time.Sleep(1 * time.Second)
		cm := exec.Command("./reload", "-pid", fmt.Sprintf("%d", cmd.Pid))
		if out, err := cm.Output(); err != nil {
			log.Println(err)
		} else {
			log.Println(strings.Trim(string(out), "\n"))
		}
	}()
	var data = make(map[string]interface{})
	data["reload"] = true
	jsonResult(c, http.StatusOK, data)
}
