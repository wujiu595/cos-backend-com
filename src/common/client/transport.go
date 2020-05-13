package client

import (
	"net"
	"net/http"
	"time"

	"github.com/wujiu2020/strip"
)

var (
	DefaultTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).Dial,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 300 * time.Second,
		ForceAttemptHTTP2:     false,
	}
)

const (
	reqIdHeader = "X-Reqid"
)

type TransportCanceler interface {
	CancelRequest(*http.Request)
}

type Transporter interface {
	http.RoundTripper
	TransportCanceler
}

type RoundTripper interface {
	Transporter
	With(http.RoundTripper) RoundTripper
}

type getReqid interface {
	ReqId() string
}

type transport struct {
	http.RoundTripper

	logger strip.Logger
}

func NewTransportWithLogger(logger strip.Logger) RoundTripper {
	return &transport{
		RoundTripper: DefaultTransport,
		logger:       logger,
	}
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	start := time.Now()
	reqId := req.Header.Get(reqIdHeader)
	if req.Header.Get(reqIdHeader) == "" {
		if log, ok := t.logger.(getReqid); ok {
			reqId = log.ReqId()
			req.Header.Set(reqIdHeader, reqId)
		}
	}
	resp, err = t.RoundTripper.RoundTrip(req)
	logTransportRequest(reqId, t.logger, req, resp, start, err)
	return
}

func (t *transport) With(rt http.RoundTripper) RoundTripper {
	newRt := *t
	newRt.RoundTripper = rt
	return &newRt
}

func (t *transport) CancelRequest(req *http.Request) {
	if rt, ok := t.RoundTripper.(TransportCanceler); ok {
		rt.CancelRequest(req)
	}
}

func logTransportRequest(reqId string, log strip.Logger, req *http.Request, resp *http.Response, start time.Time, err error) {
	var (
		respReqId string
		code      int
		extra     string

		uri      = req.URL.String()
		elaplsed = time.Since(start)
	)

	if resp != nil {
		respReqId = resp.Header.Get(reqIdHeader)
		code = resp.StatusCode
	}

	if len(respReqId) > 0 && respReqId != reqId {
		extra = ", RespReqId: " + respReqId
	}

	if err != nil {
		extra += ", Err: " + err.Error()
		if er, ok := err.(respError); ok {
			extra += ", " + er.ErrorDetail()
		}
	}

	log.Infof("Service: %s %s, Code: %d%s, Time: %dms", req.Method, uri, code, extra, elaplsed.Nanoseconds()/1e6, strip.LineOpt{Hidden: true})
}

type respError interface {
	ErrorDetail() string
}
