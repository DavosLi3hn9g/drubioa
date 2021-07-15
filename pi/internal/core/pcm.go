package core

import (
	"io"
	"os"
)

type PCM struct {
}

func (p PCM) Read(file string) ([]byte, error) {
	var (
		rate      = 16000 //采样率
		second    = 1
		frameSize = rate * 2 * second //1s的数据大小：采样率*一个采样点的数据大小2byte
		buffer    = make([]byte, frameSize)
		buf       = make([]byte, 0, frameSize)
		aWord     = make([]byte, frameSize, frameSize*3)
		isEnd     = false
	)

	audioFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	for {
		num, err := audioFile.Read(buffer)
		if err != nil {
			if err == io.EOF { //文件读取完了，改变status = STATUS_LAST_FRAME
				isEnd = true
			} else {
				panic(err)
			}
		}
		buf = buffer[:num]
		if num == 0 {
			isEnd = true
		}
		if isEnd {
			return aWord, nil
		} else {
			aWord = append(aWord, buf...)
		}

	}

}
