package audio

import (
	"VGO/pi/internal/config"
	file2 "VGO/pi/internal/pkg/file"
	"context"
	"encoding/binary"
	"github.com/youpy/go-wav"
	"io"
	"os"
	"strings"
)

var configENV = config.ENV
var logIO *file2.Log

func ReadLocal(file string, ch chan []byte, cancel context.CancelFunc) {
	var (
		buffer    = make([]byte, 640)
		buf       = make([]byte, 640)
		audioFile *os.File
		err       error
	)
	audioFile, err = os.Open(file)
	if err != nil {
		panic(err)
	}
	var isEnd bool
	for {
		num, err := audioFile.Read(buffer)
		if err != nil {
			if err == io.EOF {
				num = 0
				cancel()
				isEnd = true
			} else {
				panic(err)
			}
		}
		buf = buffer[:num]
		ch <- buf
		if num == 0 || isEnd {
			return
		}
	}

}
func Pcm2WavFile(inPath, outPath string) *os.File {
	outfile, err := os.Create(outPath)
	if err != nil {
		logIO.Fatal(err)
		return nil
	}

	audioFile, err := os.Open(inPath)
	if err != nil {
		panic(err)
	}

	var (
		frameSize    = 16000 * 2 * 3
		buffer       = make([]byte, frameSize)
		buf          = make([]byte, 0, frameSize)
		littleEndian = true
		aWord        = make([]byte, frameSize, frameSize*10)
		b16          int16
		n            int
	)

	for {
		num, err := audioFile.Read(buffer)
		if err != nil {
			if err == io.EOF { //文件读取完了，改变status = STATUS_LAST_FRAME
				num = 0
			} else {
				panic(err)
			}
		}
		buf = buffer[:num]
		if num == 0 {
			break
		} else {
			aWord = append(aWord, buf...)
		}
	}
	var numSamples uint32 = uint32(len(aWord)/2 + 1)
	var numChannels uint16 = 1
	var sampleRate uint32 = 16000
	var bitsPerSample uint16 = 16

	writer := wav.NewWriter(outfile, numSamples, numChannels, sampleRate, bitsPerSample)
	samples := make([]wav.Sample, numSamples)
	for i := 0; i < len(aWord)-1; {
		bs := make([]byte, 2)
		bs[0] = aWord[i]
		bs[1] = aWord[i+1]
		if littleEndian {
			b16 = int16(binary.LittleEndian.Uint16(bs))
		} else {
			b16 = int16(binary.BigEndian.Uint16(bs))
		}
		samples[n].Values[0] = int(b16)
		samples[n].Values[1] = 0
		n++
		i = i + 2

	}
	err = writer.WriteSamples(samples)
	if err != nil {
		logIO.Fatal(err)
		return nil
	}
	if configENV["pcm_path"] == "" {
		err = os.Remove(inPath)
		if err != nil {
			logIO.Fatal(err)
			return nil
		}
	}
	defer func() {
		outfile.Close()
	}()
	return outfile
}
func Pcm2Wav(infile string) string {
	var inPath string
	if configENV["pcm_path"] != "" {
		inPath = configENV["home_path"] + configENV["pcm_path"] + infile
	} else {
		inPath = configENV["home_path"] + "data/pcm/" + infile
	}
	outPath := configENV["home_path"] + configENV["wav_path"] + infile
	outfile := Pcm2WavFile(inPath, strings.TrimSuffix(outPath, ".pcm")+".wav")
	return outfile.Name()
}
