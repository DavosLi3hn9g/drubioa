package aliyun

import (
	"VGO/pkg/curl"
	"encoding/json"
	"errors"
	"io/ioutil"
)

var (
	ttsUrl = "https://nls-gateway.cn-shanghai.aliyuncs.com/stream/v1/tts"
)

type TTSConfig struct {
	Voice  string //发音，如Xiaowei
	Volume string //音量
	Speech string //语速，如100
	Pitch  string //语调
}
type Voice struct {
	Name  string
	Value string
	Type  string
	Note  string
}

func (isi *ClientISI) InitTTS() {
	isi.TTSConfig.Voice = setting.TTSVoice
	isi.TTSConfig.Volume = setting.TTSVolume
	isi.TTSConfig.Speech = setting.TTSSpeech
	isi.TTSConfig.Pitch = setting.TTSPitch
}

func (isi *ClientISI) TestTTS(c TTSConfig) {
	isi.TTSConfig = c
}

func (isi *ClientISI) TTS(text string, fileName string) ([]byte, error) {

	var resp RespDataISI

	if text == "" {
		return nil, errors.New("错误：text为空！")
	}
	params := map[string]string{
		"appkey":      isi.ISIappKey,
		"text":        text,
		"format":      "pcm",
		"sample_rate": "16000",
		"voice":       isi.TTSConfig.Voice,
		"volume":      isi.TTSConfig.Volume,
		"speech_rate": isi.TTSConfig.Speech,
		"pitch_rate":  isi.TTSConfig.Pitch,
	}
	bytesData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	lenght := len(bytesData)

	h := curl.Config{
		Headers: map[string]string{
			"X-NLS-Token":    isi.token,
			"Content-Type":   "application/json",
			"Content-Length": string(lenght),
		},
	}
	data, header := h.POSTJSON(ttsUrl, bytesData)
	_ = json.Unmarshal(data, &resp)

	if header != nil {
		if header.Get("Content-Type") == "audio/mpeg" {
			d1 := []byte(data)
			//  保存合成音频文件
			if fileName != "" {
				text = fileName
			}
			err = ioutil.WriteFile(configENV["home_path"]+configENV["cache_path"]+text+".pcm", d1, 0666)
			if err != nil {
				return nil, errors.New("WriteFile：" + err.Error())
			}
			return d1, nil
		}
	}
	logIO.Printf("%+q", resp)
	return nil, errors.New("音频合成失败，未知错误！")
}

var VoiceList = []Voice{
	{Name: "小云", Value: "Xiaoyun", Type: "标准女声", Note: "支持中文及中英文混合场景"},
	{Name: "小刚", Value: "Xiaogang", Type: "标准男声", Note: "支持中文及中英文混合场景"},
	{Name: "小梦", Value: "Xiaomeng", Type: "标准女声", Note: "支持中文及中英文混合场景"}, //del
	{Name: "小威", Value: "Xiaowei", Type: "标准男声", Note: "支持中文及中英文混合场景"},  //del
	{Name: "若兮", Value: "Ruoxi", Type: "温柔女声", Note: "支持中文及中英文混合场景"},
	{Name: "思琪", Value: "Siqi", Type: "温柔女声", Note: "支持中文及中英文混合场景"},
	{Name: "思佳", Value: "Sijia", Type: "标准女声", Note: "支持中文及中英文混合场景"},
	{Name: "思诚", Value: "Sicheng", Type: "标准男声", Note: "支持中文及中英文混合场景"},
	{Name: "艾琪", Value: "Aiqi", Type: "温柔女声", Note: "支持中文及中英文混合场景"},
	{Name: "艾佳", Value: "Aijia", Type: "标准女声", Note: "支持中文及中英文混合场景"},
	{Name: "艾诚", Value: "Aicheng", Type: "标准男声", Note: "支持中文及中英文混合场景"},
	{Name: "艾达", Value: "Aida", Type: "标准男声", Note: "支持中文及中英文混合场景"},
	{Name: "宁儿", Value: "Ninger", Type: "标准女声", Note: "仅支持纯中文场景"},
	{Name: "瑞琳", Value: "Ruilin", Type: "标准女声", Note: "仅支持纯中文场景"},
	{Name: "阿美", Value: "Amei", Type: "甜美女声", Note: "支持中文及中英文混合场景"},    //del
	{Name: "小雪", Value: "Xiaoxue", Type: "温柔女声", Note: "支持中文及中英文混合场景"}, //del
	{Name: "思悦", Value: "Siyue", Type: "温柔女声", Note: "支持中文及中英文混合场景"},
	{Name: "艾雅", Value: "Aiya", Type: "严厉女声", Note: "支持中文及中英文混合场景"},
	{Name: "艾夏", Value: "Aixia", Type: "亲和女声", Note: "支持中文及中英文混合场景"},
	{Name: "艾美", Value: "Aimei", Type: "甜美女声", Note: "支持中文及中英文混合场景"},
	{Name: "艾雨", Value: "Aiyu", Type: "自然女声", Note: "支持中文及中英文混合场景"},
	{Name: "艾悦", Value: "Aiyue", Type: "温柔女声", Note: "支持中文及中英文混合场景"},
	{Name: "艾婧", Value: "Aijing", Type: "严厉女声", Note: "支持中文及中英文混合场景"},
	{Name: "小美", Value: "Xiaomei", Type: "甜美女声", Note: "支持中文及中英文混合场景"},
	{Name: "艾娜", Value: "Aina", Type: "浙普女声", Note: "仅支持纯中文场景"},
	{Name: "伊娜", Value: "Yina", Type: "浙普女声", Note: "仅支持纯中文场景"},
	{Name: "思婧", Value: "Sijing", Type: "严厉女声", Note: "仅支持纯中文场景"},
	{Name: "思彤", Value: "Sitong", Type: "儿童音", Note: "仅支持纯中文场景"},
	{Name: "小北", Value: "Xiaobei", Type: "萝莉女声", Note: "仅支持纯中文场景"},
	{Name: "艾彤", Value: "Aitong", Type: "儿童音", Note: "仅支持纯中文场景"},
	{Name: "艾薇", Value: "Aiwei", Type: "萝莉女声", Note: "仅支持纯中文场景"},
	{Name: "艾宝", Value: "Aibao", Type: "萝莉女声", Note: "仅支持纯中文场景"},
	{Name: "Halen", Value: "Halen", Type: "英音女声", Note: "仅支持英文场景"}, //del
	{Name: "Harry", Value: "Harry", Type: "英音男声", Note: "仅支持英文场景"},
	{Name: "Eric", Value: "Eric", Type: "英音男声", Note: "仅支持英文场景"},
	{Name: "Emily", Value: "Emily", Type: "英音女声", Note: "仅支持英文场景"},
	{Name: "Luna", Value: "Luna", Type: "英音女声", Note: "仅支持英文场景"},
	{Name: "Luca", Value: "Luca", Type: "英音男声", Note: "仅支持英文场景"},
	{Name: "Wendy", Value: "Wendy", Type: "英音女声", Note: "仅支持英文场景"},
	{Name: "William", Value: "William", Type: "英音男声", Note: "仅支持英文场景"},
	{Name: "Olivia", Value: "Olivia", Type: "英音女声", Note: "仅支持英文场景"},
	{Name: "姗姗", Value: "Shanshan", Type: "粤语女声", Note: "支持标准粤文（简体）及粤英文混合场景"},
	{Name: "Chuangirl", Value: "chuangirl", Type: "四川话女声", Note: "中文及中英文混合场景"},
	{Name: "Lydia", Value: "lydia", Type: "英中双语女声", Note: "支持标准粤文（简体）及粤英文混合场景"},
	{Name: "艾硕", Value: "aishuo", Type: "自然男声", Note: "中文及中英文混合场景"},
	{Name: "青青", Value: "qingqing", Type: "台湾话女声", Note: "中文场景"},
	{Name: "翠姐", Value: "cuijie", Type: "东北话女声", Note: "中文场景"},
	{Name: "小泽", Value: "xiaoze", Type: "湖南重口音男声", Note: "中文场景"},
	{Name: "智香", Value: "tomoka", Type: "日语女声", Note: "日文场景"},
	{Name: "智也", Value: "tomoya", Type: "日语男声", Note: "日文场景"},
	{Name: "Annie", Value: "annie", Type: "美语女声", Note: "英文场景"},
	{Name: "佳佳", Value: "jiajia", Type: "粤语女声", Note: "标准粤文（简体）及粤英文混合场景"},
	{Name: "Indah", Value: "印尼语女声", Type: "印尼语女声", Note: "支持标准粤文及粤英文混合场景"},
	{Name: "桃子", Value: "taozi", Type: "粤语女声", Note: "支持标准粤文（简体）及粤英文混合场景"},
	{Name: "柜姐", Value: "guijie", Type: "亲切女声", Note: "支持中文及中英文混合场景"},
	{Name: "Stella", Value: "stella", Type: "知性女声", Note: "支持中文及中英文混合场景"},
	{Name: "Stanley", Value: "stanley", Type: "沉稳男声", Note: "支持中文及中英文混合场景"},
	{Name: "Kenny", Value: "kenny", Type: "沉稳男声", Note: "支持中文及中英文混合场景"},
	{Name: "Rosa", Value: "rosa", Type: "自然女声", Note: "支持中文及中英文混合场景"},
	{Name: "Farah", Value: "farah", Type: "马来语女声", Note: "仅支持纯马来语场景"},
	{Name: "马树", Value: "mashu", Type: "儿童剧男声", Note: "支持中文及中英文混合场景"},
	{Name: "小仙", Value: "xiaoxian", Type: "亲切女声", Note: "支持中文及中英文混合场景"},
	{Name: "悦儿", Value: "yuer", Type: "儿童剧女声", Note: "仅支持纯中文场景"},
}
