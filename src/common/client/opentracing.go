package client

import (
	"io"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	reqlogger "github.com/wujiu2020/strip/request-logger"
	"github.com/wujiu2020/strip/utils"
)

type openTracingTransport struct {
	http.RoundTripper
}

func NewOpenTracingTransport(tr http.RoundTripper) http.RoundTripper {
	if tr == nil {
		tr = DefaultTransport
	}
	return &openTracingTransport{tr}
}

// RoundTrip implements the RoundTripper interface.
func (t *openTracingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := t.RoundTripper

	sp, err := tracingRequestStart(req)
	if err != nil {
		return rt.RoundTrip(req)
	}

	ext.HTTPMethod.Set(sp, req.Method)
	ext.HTTPUrl.Set(sp, req.URL.String())

	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	sp.Tracer().Inject(sp.Context(), opentracing.HTTPHeaders, carrier)

	resp, err := rt.RoundTrip(req)
	if err != nil {
		sp.Finish()
		return resp, err
	}

	ext.HTTPStatusCode.Set(sp, uint16(resp.StatusCode))
	if req.Method == "HEAD" {
		sp.Finish()
	} else {
		resp.Body = closeTracker{resp.Body, sp}
	}
	return resp, nil
}

func (t *openTracingTransport) CancelRequest(req *http.Request) {
	if rt, ok := t.RoundTripper.(TransportCanceler); ok {
		rt.CancelRequest(req)
	}
}

func tracingRequestStart(req *http.Request) (sp opentracing.Span, err error) {

	var tr opentracing.Tracer
	err = utils.CtxFindValue(req.Context(), &tr)
	if err != nil {
		return
	}

	parent := opentracing.SpanFromContext(req.Context())
	var spanctx opentracing.SpanContext
	if parent != nil {
		spanctx = parent.Context()
	}

	operationName := "HTTP Client"
	root := tr.StartSpan(operationName, opentracing.ChildOf(spanctx))

	ctx := root.Context()
	sp = tr.StartSpan("HTTP "+req.Method, opentracing.ChildOf(ctx))
	ext.SpanKindRPCClient.Set(sp)

	// TODO ?? 需要增加其他 request 信息
	sp.SetTag("http.x-reqid", req.Header.Get(reqlogger.HeaderReqid))

	return
}

type closeTracker struct {
	io.ReadCloser
	sp opentracing.Span
}

func (c closeTracker) Close() error {
	err := c.ReadCloser.Close()
	c.sp.LogFields(log.String("event", "ClosedBody"))
	c.sp.Finish()
	return err
}
