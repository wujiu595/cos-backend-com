package websocket

import (
	"time"
)

const DefaultHeartbeat = 60 * time.Second

type Keepalive struct {
	Heartbeat time.Duration

	// hook close func for websocket
	Closer func()

	timer *time.Timer
}

func (c *Keepalive) Ping() string {
	c.resetTimer()
	return "pong"
}

func (c *Keepalive) Run() {
	c.resetTimer()
	<-c.timer.C

	if c.Closer != nil {
		c.Closer()
	}
	c.timer.Stop()
}

func (c *Keepalive) resetTimer() {
	if c.timer == nil {
		if c.Heartbeat <= 0 {
			c.Heartbeat = DefaultHeartbeat
		}
		c.timer = time.NewTimer(c.Heartbeat)
	}
	c.timer.Reset(c.Heartbeat)
}
