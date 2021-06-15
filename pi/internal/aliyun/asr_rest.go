package aliyun

import (
	"VGO/pkg/curl"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	audioUrl  = "https://nls-gateway.cn-shanghai.aliyuncs.com/stream/v1/asr"
	rate      = 16000 //采样率
	second    = 1
	frameSize = rate * 2 * second //1s的数据大小：采样率*一个采样点的数据大小2byte
	//intervel     = 10 * time.Second  //发送音频间隔
	//ch_reconn    = make(chan bool)
	noTalking = 0.2 //最长静默时间，用做断句
	format    = "pcm"
	endPoint  = int(2 * float64(rate) * noTalking)
)

type Chat struct {
	Key       int
	Text      string
	StartTime time.Time
	EndTime   time.Time
}

func (isi *ClientISI) ReadLoop(ch chan []byte, chStr chan *Chat) {

	var (
		buffer   = make([]byte, frameSize)
		mbufLast = make([]byte, 0, frameSize)
		mbufNext = make([]byte, 0, frameSize)
		buf      = make([]byte, 0, frameSize)
		aWord    = make([]byte, frameSize, frameSize*10)

		PauseNum     int
		isEnd        = false
		bSum, b16    int
		avg          int
		n            int
		littleEndian = true
		key          = 0
		wg           sync.WaitGroup
		wgApi        sync.WaitGroup
		chat         *Chat
	)
	wg.Add(1)
	go func() {
		for mbuf := range ch {
			if len(mbuf) == 0 {
				isEnd = true
			} else {
				if cap(buffer) >= len(buffer)+len(mbuf) {
					buffer = append(buffer, mbuf...)
					continue
				}
			}
			buf = buffer
			buffer = buffer[:0] //清空
			if len(buf) < len(mbuf) {
				logIO.Warning("注意：ReadLocal的buf容量太大！")
			}
			if chat == nil {
				chat = new(Chat)
			}
			for i := 0; i < len(buf)-1; {
				n++
				bs := make([]byte, 2)
				if littleEndian {
					bs[0] = buf[i]
					bs[1] = buf[i+1]
					b16 = int(int16(binary.LittleEndian.Uint16(bs)))
				} else {
					bs[0] = buf[i]
					bs[1] = buf[i+1]
					b16 = int(int16(binary.BigEndian.Uint16(bs)))
				}

				i = i + 2
				if 0 > b16 {
					b16 = -b16
				}
				bSum += b16

				if n >= endPoint { // 样本点结束位置，每区间0.2秒
					avg = bSum / n //average value
					n = 0
					bSum = 0
					if avg < 10 { // 发现静默，停顿段

						mbufLast = buf[:i]
						mbufNext = buf[i:]
						PauseNum++
						if PauseNum == 1 {
							if !chat.StartTime.IsZero() {
								chat.EndTime = time.Now()
							}
						} else if !isEnd {
							logIO.PrintDev(".")
						}
						break
					} else {
						if PauseNum > 0 {
							chat.StartTime = time.Now()
						}
						logIO.PrintDev("|")
						PauseNum = 0
					}
				}

			}
			if PauseNum >= 1 || isEnd {
				if PauseNum == 1 {
					aWord = append(aWord, mbufLast...)
					mbufLast = mbufLast[:0]
					logIO.PrintlnDev("+")
					wgApi.Add(1)
					if setting.PerformanceMode > 1 {
						wgApi.Done()
					}
					go func() {
						key++
						resp := isi.ASR(aWord)
						if resp.Result != "" {
							str := resp.Result
							chat.Key = key
							chat.Text = str
							chStr <- chat
							chat = nil
						} else if resp.Status == 40000005 {
							logIO.Error(fmt.Sprintf("请求太频繁，需启用高性能模式 \r\n云服务错误： %d  %s", resp.Status, resp.Message))
							//isEnd = true
						} else if resp.Status > 20000000 {
							logIO.Error(fmt.Sprintf("云服务错误： %d  %s", resp.Status, resp.Message))
							//isEnd = true
						}
						if setting.PerformanceMode <= 1 {
							wgApi.Done()
						}
					}()
					wgApi.Wait()
				}
				aWord = aWord[:0]
				aWord = append(aWord, mbufNext...)
				mbufNext = mbufNext[:0]
			} else if PauseNum == 0 {
				aWord = append(aWord, buf...)
			}

			if isEnd {
				wg.Done()
				return
			}

		}
	}()
	wg.Wait()
	close(chStr)
	logIO.Println("退出ASR实时识别！等待记录生成...")
	return
}

func (isi *ClientISI) ASR(buf []byte) RespDataISI {

	var resp RespDataISI
	length := len(buf)
	if length > 0 {

	}
	if length == 0 {
		logIO.Println("无数据！", buf)
		return resp
	}

	params := map[string]string{
		"appkey":                            isi.ISIappKey,
		"format":                            format,
		"sample_rate":                       strconv.Itoa(rate),
		"enable_punctuation_prediction":     "true",  //是否在后处理中添加标点
		"enable_inverse_text_normalization": "false", //是否在后处理中执行ITN
		"enable_voice_detection":            "false", //是否启动语音检测。说明：如果开启语音检测，服务端会对上传的音频进行静音检测，切除静音部分和之后的语音内容，不再对其进行识别；不同的模型表现结果不同。

	}

	h := curl.Config{
		Headers: map[string]string{
			"X-NLS-Token":    isi.token,
			"Content-Type":   "application/octet-stream",
			"Content-Length": string(length),
		},
	}
	data := h.POSTFILE(audioUrl, params, buf)

	err := json.Unmarshal(data, &resp)
	if err != nil {
		panic(err)
	}

	logIO.Printf("来电：%v \n", resp.Result)
	return resp
}
