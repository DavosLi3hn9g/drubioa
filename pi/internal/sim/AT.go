package sim

import (
	"VGO/pkg/fun"
	"bytes"
	"encoding/hex"
	"github.com/tarm/serial"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"
)

type CMD struct {
	Return string
	Value  string
}

type AT struct {
	Port *serial.Port
	mu   sync.Mutex
}

func (at *AT) AT(cmd string) {
	at.mu.Lock()
	if _, err := at.Port.Write([]byte(cmd + "\r\n")); err != nil {
		logIO.Fatal("AT:", err)
	}
	at.mu.Unlock()
}
func (at *AT) ATByte(cmd []byte) {
	at.mu.Lock()
	if _, err := at.Port.Write(cmd); err != nil {
		logIO.Fatal("ATByte:", err)
	}
	at.mu.Unlock()
}
func (at *AT) Call(tel string) {
	time.Sleep(time.Second)
	at.AT("ATD" + tel + ";")
}
func (at *AT) PushSMS(tel string, text string) {
	tel = strings.ToUpper(fun.HexUTF16FromString(tel))
	//fmt.Println("tel:" + tel)
	text = strings.ToUpper(fun.HexUTF16FromString(text))
	//fmt.Println("text:" + text)
	b := []byte(text)
	size := 280 //超过70个字拆分
	for i := 0; i < len(b); i += size {
		end := i + size
		if end > len(b) {
			end = len(b)
		}
		at.AT(`AT+CMGS="` + tel + `"`)
		time.Sleep(time.Second)
		at.ATByte(b[i:end])
		time.Sleep(time.Second / 5)
		b16, _ := hex.DecodeString("1A")
		at.ATByte(b16)
		time.Sleep(time.Second * 3)
	}
}
func (at *AT) ReadSMS(k string) {
	at.AT("AT+CMGR=" + k)       //获取短信
	time.Sleep(3 * time.Second) //给点时间读取和写库
	//at.AT("AT+CMGD=" + k)       //删除短信
	//time.Sleep(1 * time.Second) //不能太快的发送AT
}

func (at *AT) ReadAT() {
	var (
		str []byte
		buf = make([]byte, 8)
	)

	for {
		n, err := at.Port.Read(buf)
		if err != nil {
			logIO.Println("ReadAT", err)
			return
		}
		str = append(str, buf[:n]...)
		if (n > 0 && n < 8) || (n == 8 && bytes.HasSuffix(buf[:n], []byte{10})) {
			data := string(str)
			str = str[:0]
			logIO.Write(data)
		}
	}
}

func (at *AT) LoopWatchAT(ch chan *CMD) {
	var (
		str      []byte
		phoneNum = ""
		smsNum   = ""
		CMGL     = ""
		CMGR     = ""
	)

	for {
		var buf = make([]byte, 8)
		if at.Port == nil {
			//logIO.Println("WatchAT，待机中，正在等待重连...", at.Port)
			time.Sleep(5 * time.Second)
			continue
		}
		n, err := at.Port.Read(buf) //没有任何AT数据读取时会阻塞在这
		if err != nil {
			if err == io.EOF {
				time.Sleep(2 * time.Second)
				continue
			}
			logIO.Fatal("AT读取故障！WatchAT:", err)
			return
		}
		str = append(str, buf[:n]...)
		if bytes.HasSuffix(str, []byte{13, 10}) || bytes.HasSuffix(str, []byte{10, 10}) { //换行结尾 [13 10],注意读取SMS时的换行不一定是结尾
			data := string(str)
			data = fun.PregReplace("\r\r\n", "\r\n", data)
			str = str[:0]
			//短信必须优先处理，防止通过短信文本非法注入
			if fun.Strpos(data, "+CMGL:") > -1 || CMGL != "" { //首次发现CMGL即开始拼接，并防止因换行导致匹配失效
				CMGL += data //拼接被裁切的SMS
				if fun.Strpos(data, "OK\r\n") > -1 {
					ch <- &CMD{"CMGL", CMGL}
					CMGL = "" //OK了，结束拼接
				}
				continue
			}
			if fun.Strpos(data, "CMGR:") > -1 || CMGR != "" {
				CMGR += data
				if fun.Strpos(data, "OK\r\n") > -1 {
					ch <- &CMD{"CMGR", CMGR}
					CMGR = ""
				}
				continue
			}
			logIO.Println("----- AT -----\r\n", data)
			//新短信
			if fun.Strpos(data, "CMTI:") > -1 {
				matches := regexp.MustCompile(`\+CMTI: "(.+?)",(\d+)`).FindAllStringSubmatch(data, -1)
				if len(matches) > 0 {
					if len(matches[0]) > 1 {
						smsNum = matches[0][2]
						logIO.Printf("新短信: %v", smsNum)
						ch <- &CMD{"SMS_NUM", smsNum}
					}
				} else {
					logIO.Println("无法识别CMTI")
				}
				continue
			}
			//新来电
			if fun.Strpos(data, "RING") > -1 || fun.Strpos(data, "+CLIP:") > -1 {
				ch <- &CMD{"RING", ""}
				matches := regexp.MustCompile(`CLIP: "(.+?)",(\d+?),`).FindAllStringSubmatch(data, -1)
				if len(matches) > 0 {
					if len(matches[0]) > 1 {
						phoneNum = matches[0][1]
						logIO.Printf("来电号码: %v", phoneNum)
						ch <- &CMD{"PHONE_NUM", phoneNum}
					}
				} else {
					logIO.Println("无法识别CLIP")
				}
				continue
			}
			//通话开始
			if fun.Strpos(data, "CALL: BEGIN") > -1 {
				ch <- &CMD{"BEGIN", ""}
				continue
			}

			//对方正忙
			if fun.Strpos(data, "BUSY") > -1 {
				ch <- &CMD{"BUSY", ""}
				continue
			}
			//拒接
			if fun.Strpos(data, "MISSED_CALL") > -1 {
				ch <- &CMD{"MISSED", ""}
				continue
			}
			//短信发送无响应,可能欠费了
			if fun.Strpos(data, "CMS ERROR: Unknown error") > -1 {
				ch <- &CMD{"CMS_ERROR", ""}
				continue
			}
			//短信满了
			if fun.Strpos(data, "SMS FULL") > -1 {
				ch <- &CMD{"SMS_FULL", ""}
				continue
			}
			//通话结束
			if fun.Strpos(data, "CALL: END") > -1 {
				ch <- &CMD{"END", ""}
				continue
			}
			//未拨通或者被挂断
			if fun.Strpos(data, "NO CARRIER") > -1 {
				ch <- &CMD{"NOCARRIER", ""}
				continue
			}
			//读取CPCMFRM值，确保CPCMFRM=1
			if fun.Strpos(data, "CPCMFRM: 1") > -1 {
				ch <- &CMD{"CPCMFRM", "1"}
				continue
			} else if fun.Strpos(data, "CPCMFRM: 0") > -1 {
				ch <- &CMD{"CPCMFRM", "0"}
				continue
			} else if fun.Strpos(data, "CPCMFRM: ") > -1 {
				ch <- &CMD{"CPCMFRM", ""}
				continue
			}
			if fun.Strpos(data, "OK") == 0 {
				ch <- &CMD{"OK", ""}
				continue
			}

		}

	}
}
