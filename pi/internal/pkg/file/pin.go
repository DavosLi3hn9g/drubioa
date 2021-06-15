package file

import (
	"VGO/pi/internal/pkg/logfile"
	"encoding/json"
)

type Pin struct {
	Pin   uint8
	State uint8
}

func (l *Log) PinState(pin, state uint8) {
	var m Pin
	m.Pin = pin
	m.State = state
	b, _ := json.Marshal(&Msg{
		T: 2,
		M: m,
	})
	if logIO.WsOut != nil {
		_, err := logIO.WsOut.Write(b)
		if err != nil {
			logfile.Error(err)
		}
	}

}
