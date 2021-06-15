package cache

import (
	"VGO/pi/internal/orm"
	"VGO/pkg/fun"
	"strconv"
	"strings"
)

var SettingsCache = new(Settings)

type Settings struct {
	NotNil          bool
	TalkPrologue    string
	TTSVoice        string
	TTSSpeech       string
	TTSPitch        string
	TTSVolume       string
	TTYCall         string
	TTYUSB          string
	AliyunAkId      string
	AliyunAkSecret  string
	AliyunIsiAppkey string
	AliyunBucket    AliyunBucket
	UploadOSS       bool
	UnitPrice       float64
	PerformanceMode int
}
type AliyunBucket struct {
	Location string
	Name     string
}

func (s Settings) Read() *Settings {
	if !SettingsCache.NotNil {
		var set = make(map[string]string)
		var uploadOSS bool
		var unitPrice float64
		var performanceMode int
		all := orm.Settings{}.All()
		for _, v := range all {
			set[v.Key] = v.Value
		}
		if set["upload_oss"] != "" {
			uploadOSS = true
		}
		if set["unit_price"] != "" {
			float, _ := strconv.ParseFloat(set["unit_price"], 32)
			unitPrice = fun.Round(float, 2)
		}
		if set["performance_mode"] != "" {
			performanceMode, _ = strconv.Atoi(set["performance_mode"])
		}
		*SettingsCache = Settings{
			NotNil:          true,
			TalkPrologue:    set["talk_prologue"],
			TTSVoice:        set["tts_voice"],
			TTSSpeech:       set["tts_speech"],
			TTSPitch:        set["tts_pitch"],
			TTSVolume:       set["tts_volume"],
			TTYCall:         set["tty_call"],
			TTYUSB:          set["tty_usb"],
			AliyunAkId:      set["aliyun_ak_id"],
			AliyunAkSecret:  set["aliyun_ak_secret"],
			AliyunIsiAppkey: set["aliyun_isi_appkey"],
			AliyunBucket:    AliyunBucket{},
			UploadOSS:       uploadOSS,
			UnitPrice:       unitPrice,
			PerformanceMode: performanceMode,
		}
		if set["aliyun_bucket_name"] != "" {
			bucket := strings.Split(set["aliyun_bucket_name"], ".")
			if len(bucket) > 1 {
				SettingsCache.AliyunBucket.Name = bucket[0]
				SettingsCache.AliyunBucket.Location = bucket[1]
			}
		}

	}
	return SettingsCache
}

func (s Settings) Update() *Settings {
	s.Clear()
	return s.Read()
}

func (s Settings) Clear() {
	*SettingsCache = Settings{}
}
