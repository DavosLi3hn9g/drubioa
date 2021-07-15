package core

import (
	"VGO/pi/internal/aichat"
	"VGO/pi/internal/aliyun"
	"VGO/pi/internal/audio"
	"VGO/pi/internal/cache"
	"VGO/pi/internal/cmd"
	"VGO/pi/internal/config"
	"VGO/pi/internal/orm"
	"VGO/pi/internal/pkg/file"
	"VGO/pi/internal/sim"
	"context"
	"encoding/json"
	"github.com/tarm/serial"
	"strconv"
	"strings"
	"time"
)

type Talk struct {
	IsCaller  bool          `json:"is_caller"`  //是否主叫方
	BeginTime time.Duration `json:"begin_time"` //结束时间偏移
	EndTime   time.Duration `json:"end_time"`   //结束时间偏移
	Text      string        `json:"text"`       //内容
}

var (
	logIO       *file.Log
	SerialAT    = new(sim.AT)
	SerialAudio = new(sim.Audio)
	isi         = new(aliyun.ClientISI)
	PhoneCore   = new(Phone)
	Qiar        = new(aichat.IQiar)
	configENV   = config.ENV
	FlagReply   *string
	FlagTest    *string
	Setting     *cache.Settings
	Ctx         context.Context
	IsRunning   = new(bool)
	IsTalking   = new(bool)
	CPCMFRM     = new(bool)
	SmsList     = make(map[string]*orm.LogsSms)
)

func Start() {

	var (
		telTest   = false
		chGPIO    = make(chan string)
		outFile   = make(chan string)
		flagReply = *FlagReply
		flagTest  = *FlagTest
		err       error
		simState  string
	)

	logIO.Println("助手已启动...")
	if flagReply != "" {
		logIO.Println("自动回复：", flagReply)
	}
	if flagTest == "10010" || flagTest == "10086" || flagTest == "10000" {
		telTest = true
		go func() { chGPIO <- "START" }()
	}
	bcmRing, _ := strconv.Atoi(configENV["bcm_ring"])
	go LoopPin(uint8(bcmRing), chGPIO)
	go func() {
		chGPIO <- "INIT"
	}()

	Ctx = context.TODO()
	go PhoneCore.LoopPhoneATRead()
	for {
		select {
		case <-PhoneCore.LoopCheckAnswered(0):
			simState = "START"
			logIO.Warning("响铃故障！")
		case simState = <-chGPIO:
		}
		if simState == "START" || simState == "INIT" {
			*IsRunning = true
			go cmd.CheckNetwork()
			go func() {
				err := isi.CreateToken(nil)
				if err != nil {
					logIO.Fatal(err)
				}
			}()
			CallEnd(0, true)
			PhoneCore.Ctx, PhoneCore.Cancel = context.WithCancel(Ctx)
			if SerialAT.Port == nil {
				SerialAT.Port, err = serial.OpenPort(&serial.Config{Name: Setting.TTYCall, Baud: 115200})
				if err != nil {
					if err.Error() == "The requested resource is in use." {
						logIO.Warning("" + Setting.TTYCall + " 串口正在被占用，请检查后重启本应用。")
					} else {
						logIO.Fatal(err, "请检查"+Setting.TTYCall+"串口连接是否正常？")
					}
					continue
				}
			}
			if simState == "INIT" {
				err = SerialAT.Port.Flush()
				if err != nil {
					logIO.Warning("Flush: ", err)
					return
				}
				PhoneCore.PiInitAT()
				*IsRunning = false
				CallEnd(10, false)
				continue
			}
			if SerialAudio.Port == nil {

				SerialAudio.Port, err = serial.OpenPort(&serial.Config{Name: Setting.TTYUSB, Baud: 115200, ReadTimeout: time.Second * 20})
				if err != nil {
					logIO.Fatal(err, "请检查"+Setting.TTYUSB+" USB串口连接是否正常？")
					time.Sleep(10 * time.Second)
					if *IsTalking {
						SerialAT.AT("AT+CHUP")
					}
					continue
				}
			}
			//log.Printf("Setting：%+v", Setting)
			logIO.Println("串口准备就绪...")
			if telTest {
				logIO.Println("测试电话语音识别：拨号", flagTest)
				SerialAT.AT("ATD" + flagTest + ";")
			}
			var (
				chASR     = make(chan *aliyun.Chat, 10)
				outTTS    = make(chan string, 10)
				record    = make([]Talk, 0, 20)
				rec       *aliyun.Record
				smsText   string
				smsJson   string
				startTime time.Time
				beginTime time.Duration
				pcmFile   string
			)
			select {
			case <-PhoneCore.LoopCheckAnswered(0):
				logIO.Println("开始通话！")
				go func() { outTTS <- Setting.TalkPrologue }() //"你好，有什么事？"
				go PhoneCore.PhoneAudioRead(chASR, outFile)
				if LogCall.TimeStart.Unix() > 0 {
					startTime = LogCall.TimeStart
				} else {
					startTime = time.Now()
				}
				beginTime = time.Since(startTime)
				Qiar.IsEND = false
				go func() {
					for inChat := range chASR {
						if inChat.Text != "" {
							record = append(record, Talk{IsCaller: true, BeginTime: beginTime, EndTime: time.Since(startTime), Text: inChat.Text})
							outTTS <- Qiar.Out(inChat.Text, flagReply)
						}
						if Qiar.IsEND {
							go func() {
								time.Sleep(20 * time.Second) //给个20秒说拜拜！
								SerialAT.AT("AT+CHUP")       //PhoneCore.Cancel()
								logIO.Warning("已触发挂断策略！")
							}()
							Qiar.IsEND = false
						}
					}
					close(outTTS)
				}()
				go func() {
					for out := range outTTS {
						if out != "" {
							beginTime = time.Since(startTime)
							PhoneCore.PhoneAudioWrite(out)
							record = append(record, Talk{IsCaller: false, BeginTime: beginTime, EndTime: time.Since(startTime), Text: out})
						}
					}
				}()
				select {
				case <-PhoneCore.Ctx.Done():
				case <-time.After(time.Second * 120):
					SerialAT.AT("AT+CHUP")
					logIO.Warning("通话时间超过2分钟，强制挂断！")
				}
				select {
				case pcmFile = <-outFile:
					talkList := ""
					for _, r := range record {
						if r.IsCaller {
							talkList += r.Text
						}
					}
					smsText = talkList
					smsByte, _ := json.Marshal(record)
					smsJson = string(smsByte)
					if smsText != "" {
						LogCall.End(smsText, smsJson, pcmFile)
						PhoneCore.Ctx, PhoneCore.Cancel = context.WithTimeout(Ctx, 60*time.Second) //超60s强制终止拨打
						if Qiar.User.Tel != "" && Qiar.Policies.Checked != "" {
							checked := orm.Policies{}.CheckedToStruct(Qiar.Policies.Checked)
							if checked.Sms {
								logIO.Println("发送短信：", smsText)
								SerialAT.PushSMS(Qiar.User.Tel, PhoneCore.NmbFrom+"来电："+smsText)
							}
							if checked.Call {
								*IsRunning = true
								logIO.Println("拨打电话：", Qiar.User.Tel)
								SerialAT.Call(Qiar.User.Tel)
								select {
								case <-PhoneCore.LoopCheckAnswered(0):
									PhoneCore.PhoneAudioWrite("有重要来电，注意查看短信！")
								case <-time.After(time.Second * 60):
									logIO.Println("拨出电话超60秒无人接听！")
								}
								<-PhoneCore.Ctx.Done() //等待用户挂断或超时强制挂断
								if *IsRunning {
									SerialAT.AT("AT+CHUP")
								}
							}
							Qiar.Policies = nil
						}
					}
					wavName := audio.Pcm2Wav(pcmFile + ".pcm")
					if Setting.UploadOSS {
						oss := new(aliyun.ClientOSS).NewClient(nil)
						ossName := strings.TrimPrefix(wavName, "./")
						oss.Upload("call/"+ossName, wavName, Setting.AliyunBucket.Name)
						logIO.Println("录音已上传：", wavName)
						pcmURL := oss.GetURL("call/"+ossName, Setting.AliyunBucket.Name)
						isi.InitRecord()
						rec = isi.Record(pcmURL, configENV["home_path"]+configENV["wav_path"]+pcmFile+".json")
						smsText = rec.Text
						smsJson = rec.Content
						LogCall.EndUpdate(smsText, smsJson, pcmFile)
					}
				case <-time.After(time.Second * 60):
					SerialAT.AT("AT+CHUP")
					logIO.Fatal("60秒没有收到文件信息，请检查端口" + Setting.TTYUSB + "读取是否正常！")
				}
			case <-time.After(time.Second * 60):
				logIO.Fatal("60秒没有收到接听信号，请检查AT读取是否正常！")
			}
			*IsTalking = false
			PhoneCore.NmbFrom = ""
			CallEnd(10, false)
		}
		simState = ""
	}

}

func CallEnd(waitSecond int, closeAT bool) {
	time.Sleep(time.Duration(waitSecond) * time.Second)
	if !*IsRunning {
		logIO.Println("CallEnd！")
		if SerialAT.Port != nil && closeAT {
			_ = SerialAT.Port.Flush()
			_ = SerialAT.Port.Close()
			SerialAT.Port = nil
		}
		if SerialAudio.Port != nil {
			_ = SerialAudio.Port.Flush()
			_ = SerialAudio.Port.Close()
			SerialAudio.Port = nil
		}
		Qiar = new(aichat.IQiar)
		logIO.Println("进入待机状态！")
		return
	}
}
