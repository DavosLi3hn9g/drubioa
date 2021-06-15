package socket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func (c *Client) WsFile(context *gin.Context) {
	c.newClient(context)

	var lastMod time.Time
	if n, err := strconv.Atoi(context.PostForm("lastMod")); err == nil {
		lastMod = time.Unix(int64(n), 0)
	}
	go c.writer(lastMod)
	go c.reader()
}

func (c *Client) writer(lastMod time.Time) {
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod)
	fileTicker := time.NewTicker(filePeriod)
	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		//文件读取时间间隔
		case <-fileTicker.C:
			var p []byte
			var err error

			p, lastMod, err = ReadFileIfModified(lastMod)

			if err != nil {
				if s := err.Error(); s != lastError {
					lastError = s
					p = []byte(lastError)
				}
			} else {
				lastError = ""
			}

			if p != nil {
				_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.conn.WriteMessage(websocket.TextMessage, p); err != nil {
					return
				}
			}
		//心跳时间间隔
		case <-pingTicker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Client) reader() {
	defer func() {
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func ReadFileIfModified(lastMod time.Time) ([]byte, time.Time, error) {
	filename := configENV["home_path"] + configENV["log_path"] + "pi.log"
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, lastMod, err
	}
	if !fi.ModTime().After(lastMod) {
		return nil, lastMod, nil
	}
	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fi.ModTime(), err
	}
	return p, fi.ModTime(), nil
}
