package websocket

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/wujiu2020/strip"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

type payload struct {
	Id     string      `json:"id,omitempty"`
	Action string      `json:"action,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func (r *payload) BindData(v interface{}) error {
	body, typ, _ := websocket.JSON.Marshal(r.Data)
	err := websocket.JSON.Unmarshal(body, typ, v)
	return err
}

type webSocketRPC struct {
	log  strip.Logger
	conn *websocket.Conn

	closed bool

	services map[string]reflect.Value
}

func (c *webSocketRPC) Handle() {
	for {
		if c.closed {
			return
		}
		var res payload
		err := websocket.JSON.Receive(c.conn, &res)
		if err == io.EOF || isConnectClosed(err) {
			return
		}
		if err != nil {
			c.log.Info("", zap.Error(err))
		}

		go c.Receive(&res)
	}
}

func (c *webSocketRPC) Close() {
	c.closed = true
}

func (c *webSocketRPC) Register(name string, service interface{}) {
	if c.services == nil {
		c.services = make(map[string]reflect.Value)
	}

	val := reflect.ValueOf(service)
	if val.Kind() != reflect.Ptr && val.Elem().Kind() != reflect.Struct {
		panic("websocket rpc service need pointer of struct")
	}

	c.services[name] = val.Elem()
}

func (c *webSocketRPC) Send(action string, data interface{}) {
	err := websocket.JSON.Send(c.conn, &payload{
		Action: action,
		Data:   data,
	})
	if isConnectClosed(err) {
		return
	}
	if err != nil {
		c.log.Warn("", zap.Error(err))
	}
}

func (c *webSocketRPC) Reply(res *payload) {
	err := websocket.JSON.Send(c.conn, res)
	if isConnectClosed(err) {
		return
	}
	if err != nil {
		c.log.Warn("", zap.Error(err))
	}
}

func (c *webSocketRPC) Receive(req *payload) {
	defer func() {
		if err := recover(); err != nil {
			c.log.Warn(fmt.Sprint(err))
		}
	}()

	if req.Action == "" {
		return
	}

	parts := strings.Split(req.Action, ".")

	var (
		name       string
		action     string
		methodFunc reflect.Value
	)

	if len(parts) == 2 {
		name, action = parts[0], parts[1]
	} else {
		name, action = "", parts[0]
	}

	action = strings.Title(action)

	elm, _ := c.services[name]
	if elm.Kind() != reflect.Invalid {
		methodFunc = elm.Addr().MethodByName(action)
	}

	if methodFunc.Kind() == reflect.Invalid {
		c.log.Warn("wsrpc wrong method", zap.String("reqAction", req.Action), zap.String("action", action), zap.String("name", name))
		return
	}

	mTyp := methodFunc.Type()
	if mTyp.NumIn() > 1 {
		c.log.Warn("wsrpc func must only have 1 in data")
		return
	}

	var in []reflect.Value

	if mTyp.NumIn() > 0 {
		argIn := mTyp.In(0)
		arg := reflect.New(argIn)
		v := arg.Interface()
		err := req.BindData(v)
		if err != nil {
			c.log.Warn("wsrpc wrong data unmarshal", zap.Error(err))
			return
		}
		in = append(in, arg.Elem())
	}

	outs := methodFunc.Call(in)
	if len(outs) == 0 {
		return
	}

	out := outs[0]
	if !out.CanInterface() {
		return
	}

	res := &payload{
		Id:   req.Id,
		Data: out.Interface(),
	}
	c.Reply(res)
}
