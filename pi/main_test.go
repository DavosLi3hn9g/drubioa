package main

import (
	"VGO/pi/internal/aliyun"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_Record(t *testing.T) {
	oss := new(aliyun.ClientOSS).NewClient(nil)
	pcmURL := oss.GetURL("data/pcm/tonghua_20190805_22_55.wav", "iqiar-pcm")
	isi := new(aliyun.ClientISI)
	isi.InitRecord()
	text := isi.Record(pcmURL, "./data/pcm/tonghua_20190804_10_46.json")
	t.Log(text)
}

func Test_TTS(t *testing.T) {

	isi := new(aliyun.ClientISI)
	isi.CreateToken(nil)
	isi.InitTTS()
	isi.TTS("反电销骚扰，个人版AI解决方案", "解说")

}

func Test_ReloadTTS(t *testing.T) {

	var (
		files  []string
		dirs   []string
		dirPth = "./cache"
	)

	isi := new(aliyun.ClientISI)
	isi.CreateToken(nil)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		t.Error(err)
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {

		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
		} else {
			ok := strings.HasSuffix(fi.Name(), ".pcm") // 指定格式
			if ok {
				t.Log(fi.Name())
				isi.InitTTS()
				isi.TTS(strings.TrimSuffix(fi.Name(), ".pcm"), "")
				files = append(files, dirPth+PthSep+fi.Name())

			}
		}
	}
}

func Test_SMS(t *testing.T) {

}
