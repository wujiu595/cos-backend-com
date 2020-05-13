package websocket

import (
	"net/http"
	"strings"

	"github.com/wujiu2020/strip"

	"golang.org/x/net/websocket"
)

type WebSocket struct {
	ws   websocket.Server
	conn *websocket.Conn
	log  strip.Logger

	rpc webSocketRPC
}

func NewWebSocket(log strip.Logger) *WebSocket {
	ws := new(WebSocket)
	ws.log = log
	return ws
}

func (c *WebSocket) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c.ws.Handler = c.handle
	c.ws.ServeHTTP(w, req)
}

func (c *WebSocket) Register(name string, v interface{}) {
	c.rpc.Register(name, v)
}

func (c *WebSocket) Send(action string, data interface{}) {
	c.rpc.Send(action, data)
}

func (c *WebSocket) Close() {
	if c.conn != nil {
		c.conn.Close()
	}

	c.rpc.Close()
}

func (c *WebSocket) handle(conn *websocket.Conn) {
	// saved conn
	c.conn = conn

	// set rpc
	c.rpc.conn = conn
	c.rpc.log = c.log

	c.rpc.Handle()
}

// https://github.com/golang/go/issues/4373
func isConnectClosed(err error) (yes bool) {
	if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
		yes = true
	}
	return
}
