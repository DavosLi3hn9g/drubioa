package socket

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type GPin struct {
	Outside []struct {
		BCM string
	}
	Inside []struct {
		BCM string
	}
}

func (c *Client) pumpStdin(w io.Writer) {
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

func (c *Client) pumpStdout(r io.Reader, done chan struct{}) {
	var ws = c.conn
	s := bufio.NewScanner(r)
	for s.Scan() {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
			ws.Close()
			break
		}
	}
	if s.Err() != nil {
		log.Println("scan:", s.Err())
	}
	close(done)
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}
func (c *Client) WsCmd(context *gin.Context) {
	return
}
