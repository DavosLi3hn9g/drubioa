package audio

import (
	"VGO/pi/internal/aichat"
	"VGO/pi/internal/aliyun"
	"context"
)

type Audio struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	ISI    *aliyun.ClientISI
}

func (p *Audio) FileRead(outASR chan *aliyun.Chat, inFile string) {
	var chAudio = make(chan []byte)
	go ReadLocal(inFile, chAudio, p.Cancel) //从本地文件读取Audio
	go p.ISI.ReadLoop(chAudio, outASR)      //识别Audio
	<-p.Ctx.Done()
	return
}

func FileASR(filePath string, reply string, ctx context.Context, cancel context.CancelFunc) context.Context {

	var outASR = make(chan *aliyun.Chat, 10)
	var outTTS = make(chan string, 10)
	var localTest = false
	var Audio = new(Audio)
	var Qiar = new(aichat.IQiar)
	Audio.Ctx, Audio.Cancel = context.WithCancel(ctx)
	if filePath != "" {
		logIO.Println("测试本地文件：", filePath)
		localTest = true
	}

	if localTest {
		Audio.ISI = new(aliyun.ClientISI)
		err := Audio.ISI.CreateToken(nil)
		if err != nil {
			logIO.Fatal(err)
		}
		go Audio.FileRead(outASR, filePath)
		go func() {
			for outChat := range outASR {
				if outChat.Text != "" {
					outTTS <- Qiar.Out(outChat.Text, reply)
				}
			}
			close(outTTS)
		}()
		go func() {
			for out := range outTTS {
				if out != "" {
					logIO.Println("AI：", out)
				}
			}
		}()
		select {
		case <-Audio.Ctx.Done():
			logIO.Println(reply + "读取完毕！")
		}
	}
	cancel()
	return Audio.Ctx
}
