package socket

import (
	"VGO/pi/internal/config"
	"VGO/pi/internal/pkg/logfile"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Client struct {
	// websocket连接
	conn *websocket.Conn
	// 发送消息的通道
	send chan []byte
}

var (
	upGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	configENV = config.ENV
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 1 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

func (c *Client) newClient(context *gin.Context) {
	w := context.Writer
	r := context.Request
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			logfile.Error(err)
		}
	}
	c.conn = conn
	//return &Client{conn, make(chan []byte, 256)}
}

func (c *Client) ping(done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Printf("ping: %v", err)
				return
			}
		case <-done:
			return
		}
	}
}

func (c *Client) writeError(msg string, err error) {
	log.Println(msg, err)
	_ = c.conn.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}
