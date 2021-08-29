package core

import (
	"VGO/pi/internal/aliyun"
	"VGO/pi/internal/cache"
	"VGO/pi/internal/orm"
	"VGO/pi/internal/pkg/file"
	"VGO/pi/internal/sim"
	"VGO/pkg/fun"
	"context"
	"encoding/hex"
	"regexp"
	"sync"
	"time"
)

type Phone struct {
	Ctx      context.Context
	Cancel   context.CancelFunc
	NmbFrom  string //主叫号码
	NmbTo    string //被叫号码
	IO       *file.Log
	isEnd    bool
	Callback map[string]string
	Chan     chan string
}

func (p *Phone) PiInitAT() {
	wait := 2 * time.Second
	SerialAT.AT("AT+CPCMFRM?") //如果是16k，返回CPCMFRM=1
	time.Sleep(wait * 3)       //给点时间拿到结果
	if !*CPCMFRM {
		SerialAT.AT("AT+CLIP=1")    //来电显示
		time.Sleep(wait)            //这里必须要停顿一下，否则会出错
		SerialAT.AT("AT+CPCMFRM=1") //16k
		time.Sleep(wait)
		SerialAT.AT("AT+CLVL=5") //最大音量
		time.Sleep(wait)
		//识别中文短信
		SerialAT.AT("AT+CMGF=1") //设置文本显示
		time.Sleep(wait)
		SerialAT.AT("AT+CSCS=\"UCS2\"") //设置GSM编码集
		time.Sleep(wait)
		SerialAT.AT("AT+CSMP=17,167,2,25") //设置文本模式参数
		time.Sleep(wait)
		SerialAT.AT("AT+CNMI=2,1") //设置新信息提醒
		time.Sleep(wait)
		b16, _ := hex.DecodeString("1B") //防止当前在短信发送状态
		SerialAT.ATByte(b16)
		time.Sleep(wait)
	}
	SerialAT.AT("AT+CMGL=\"ALL\"")
}

func (p *Phone) PhoneAudioRead(outASR chan *aliyun.Chat, outFile chan string) {
	var chAudio = make(chan []byte)
	go SerialAudio.Read(chAudio, outFile, Setting.TTYUSB) //从串口读取Audio
	go isi.ReadLoop(chAudio, outASR)                      //识别Audio
	<-p.Ctx.Done()
	return
}
func (p *Phone) PhoneAudioWrite(out string) {
	logIO.Println("AI：", out)
	b, err := PCM{}.Read(configENV["home_path"] + configENV["cache_path"] + out + ".pcm")
	if err != nil {
		logIO.Println("找不到语音缓存文件，正在云端合成，请稍等...")
		isi.InitTTS()
		b, err = isi.TTS(out, "")
		if err != nil {
			logIO.Error(err)
		}
	}
	SerialAudio.Write(b)
}
func (p *Phone) LoopPhoneATRead() {
	p.Callback = make(map[string]string)
	var (
		ch         = make(chan *sim.CMD, 10)
		prevReturn = "*"
		text       string
		timeLoc, _ = time.LoadLocation("Asia/Shanghai") //设置时区
		mu         sync.Mutex
	)
	go SerialAT.LoopWatchAT(ch)

	for at := range ch {
		mu.Lock()
		p.isEnd = false

		if at.Return != "OK" {
			logIO.PrintlnDev("\r\n+++++\r\nReturn：", at.Return, "\r\n+++++\r\n")
		}
		if at.Return == "CPCMFRM" {
			if at.Value == "1" {
				*CPCMFRM = true
			} else if at.Value == "0" {
				*CPCMFRM = false
			} else {
				*CPCMFRM = true
				logIO.Warning("当前扩展板存在异常，建议在WEB控制台执行AT+CRESET复位扩展板后重启QiarAI！")
			}
		}
		if at.Return == "SMS_NUM" {
			smsNum := at.Value
			SerialAT.AT("AT+CMGR=" + smsNum)
			logIO.Println("有新短信！！！！")
		}
		if at.Return == "RING" {
			logIO.Println("正在响铃...")
			*IsRunning = true
		}
		if at.Return == "PHONE_NUM" {
			p.NmbFrom = at.Value
			SerialAT.AT("ATA")
			logIO.Println("有新来电！！！！")
			go LogCall.Start(p.NmbFrom)
		}

		if at.Return == "BEGIN" {
			logIO.Println("已接听！")
			time.Sleep(time.Second / 5) //这里必须要停顿一下，否则会出错
			SerialAT.AT("AT+CPCMREG=1")
			time.Sleep(time.Second / 5) //这里必须要停顿一下，否则会出错
			*IsTalking = true
		}

		if at.Return == "END" {
			logIO.Println("已挂断！")
			p.isEnd = true
		}
		if at.Return == "CMS_ERROR" {
			logIO.Fatal("号码可能欠费了！")
			p.isEnd = true
		}
		if at.Return == "SMS_FULL" {
			logIO.Fatal("短信满了！请删除一些过期短信")
		}
		if at.Return == "NOCARRIER" || at.Return == "MISSED" {
			logIO.Println("终止通话！")
			p.isEnd = true
			*IsRunning = false
		}
		if at.Return == "OK" {
			p.Callback[prevReturn] = "OK"
		} else {
			p.Callback[at.Return] = ""
			prevReturn = at.Return
		}
		if p.isEnd {
			SerialAT.AT("AT+CPCMREG=0")
			*IsTalking = false
			time.AfterFunc(500*time.Millisecond, func() { //AT响应有时候需要一点时间
				p.Cancel()
			})
		}
		if at.Return == "CMGL" {
			matches := regexp.MustCompile(`CMGL: (\d+?),"(.+?)","(\w+?)","(.*?)","(.+?)"\r\n(\w+)`).FindAllStringSubmatch(at.Value, -1)
			for _, v := range matches {
				text, _ = fun.Unicode2String(v[6])
				telFrom, _ := fun.Unicode2String(v[3])
				timeGo, _ := time.ParseInLocation("06/01/02,15:04:05-07", v[5], timeLoc)
				SmsList[v[1]] = &orm.LogsSms{
					Text:     text,
					TelFrom:  telFrom,
					TelTo:    p.NmbTo,
					Dateline: int(timeGo.Unix()),
					SmsId:    v[1],
				}
			}
			*IsRunning = false
		}
		if at.Return == "CMGR" {
			matches := regexp.MustCompile(`AT\+CMGR=(\d+?)\r\n\+CMGR: "(.+?)","(\w+?)","(.*?)","(.+?)"\r\n(\w+)\r\n`).FindAllStringSubmatch(at.Value, -1)
			if len(matches) > 0 {
				if len(matches[0]) > 6 {
					id := matches[0][1]
					text, _ = fun.Unicode2String(matches[0][6])
					telFrom, _ := fun.Unicode2String(matches[0][3])
					timeGo, _ := time.ParseInLocation("06/01/02,15:04:05-07", matches[0][5], timeLoc)
					SmsList[id] = &orm.LogsSms{
						Text:     text,
						TelFrom:  telFrom,
						TelTo:    p.NmbTo,
						Dateline: int(timeGo.Unix()),
						SmsId:    id,
					}
					SerialAT.PushSMS(cache.UsersCache.Default().Tel, "来自"+telFrom+"的新短信："+text)
				}
			} else {
				logIO.Println("无法识别CMGR")
			}
		}
		mu.Unlock()
	}
}

func (p *Phone) LoopCheckAnswered(over int) <-chan string {
	var chCmd = make(chan string)

	go func() {
		var i = 0
		for {
			i++
			if *IsTalking {
				chCmd <- "BEGIN"
				return
			}
			time.Sleep(2 * time.Second)
			if i > over && over > 0 {
				return
			}
		}
	}()
	return chCmd
}
