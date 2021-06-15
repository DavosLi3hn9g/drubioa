package socket

import (
	"VGO/pi/internal/core"
	"VGO/pi/internal/pkg/file"
	"VGO/pi/internal/pkg/logfile"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"time"
)

func (c *Client) printStdin(w *io.PipeWriter) {
	var ws = c.conn
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		message = append(message, '\n')
		if _, err := w.Write(message); err != nil {
			break
		}
	}
}

func (c *Client) printStdout(reader *io.PipeReader, done chan struct{}) {

	var ws = c.conn
	buf := make([]byte, 1024) // 单次推送的最多字节数
	for {
		n, err := reader.Read(buf)
		if err != nil {
			logfile.Error(err)
		}
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
			logfile.Error(err)
			ws.Close()
			break
		}
	}
	close(done)
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
	core.PinCacheClear()
}

func (c *Client) WsPrint(context *gin.Context) {

	c.newClient(context)
	defer c.conn.Close()

	outr, inw := io.Pipe()

	defer inw.Close()
	defer outr.Close()
	//defer core.CallEnd(60, false) //发送AT命令时会读取串口，这里可以断开WS连接后终止串口读取
	if inw != nil {
		file.WsWrite(inw)
		defer file.WsClose()
	}

	stdoutDone := make(chan struct{})
	go c.printStdout(outr, stdoutDone)
	go c.printStdin(inw)
	go c.ping(stdoutDone)
	<-stdoutDone
}
