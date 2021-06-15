package socket

import (
	"bufio"
	"flag"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gorilla/websocket"
)

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
	//二进制文件路径
	cmdPath, err := exec.LookPath("sh")
	if err != nil {
		log.Fatal(err)
	}
	c.newClient(context)
	var ws = c.conn
	defer ws.Close()

	outr, outw, err := os.Pipe()
	if err != nil {
		c.writeError("stdout:", err)
		return
	}
	defer outr.Close()
	defer outw.Close()

	inr, inw, err := os.Pipe()
	if err != nil {
		c.writeError("stdin:", err)
		return
	}
	defer inr.Close()
	defer inw.Close()
	//defer core.CallEnd(60, false) //发送AT命令时会读取串口，这里可以断开WS连接后终止串口读取
	proc, err := os.StartProcess(cmdPath, flag.Args(), &os.ProcAttr{
		Files: []*os.File{inr, outw, outw},
	})
	if err != nil {
		c.writeError("start:", err)
		return
	}

	inr.Close()
	outw.Close()

	stdoutDone := make(chan struct{})
	go c.pumpStdout(outr, stdoutDone)
	go c.ping(stdoutDone)

	c.pumpStdin(inw)

	// Some commands will exit when stdin is closed.
	inw.Close()

	// Other commands need a bonk on the head.
	if err := proc.Signal(os.Interrupt); err != nil {
		log.Println("inter:", err)
	}

	select {
	case <-stdoutDone:
	case <-time.After(time.Second):
		// A bigger bonk on the head.
		if err := proc.Signal(os.Kill); err != nil {
			log.Println("term:", err)
		}
		<-stdoutDone
	}

	if _, err := proc.Wait(); err != nil {
		log.Println("wait:", err)
	}
}
