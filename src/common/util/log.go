package util

import (
	"net/http"

	reqlogger "github.com/wujiu2020/strip/request-logger"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	"github.com/wujiu2020/strip"
)

func DefaultLoggerOption(sp *strip.Strip) reqlogger.LoggerOption {
	return reqlogger.LoggerOption{
		ColorMode:     sp.Config.RunMode.IsDev(),
		LineInfo:      true,
		ShortLine:     sp.Config.RunMode.IsProd(),
		FlatLine:      sp.Config.RunMode.IsProd(),
		LogStackLevel: strip.LevelCritical,
		ReqidFilter: func(ctx strip.Context, rw http.ResponseWriter, req *http.Request, content string) string {
			sp := opentracing.SpanFromContext(ctx)
			if sp != nil {
				if span, ok := sp.Context().(jaeger.SpanContext); ok && span.SpanID() > 0 {
					return span.SpanID().String()
				}
			}
			return content
		},
		PrefixFilter: func(ctx strip.Context, rw http.ResponseWriter, req *http.Request, content string) string {
			if HasVersion() {
				return "[" + Version + "]" + content
			} else {
				return content
			}
		},
		ReqBegFilter: func(ctx strip.Context, rw http.ResponseWriter, req *http.Request, content string) string {
			if req.URL.Path == "/_version" {
				return ""
			}
			return content
		},
		ReqEndFilter: func(ctx strip.Context, rw http.ResponseWriter, req *http.Request, content string) string {
			if req.URL.Path == "/_version" {
				return ""
			}
			return content
		},
	}
}
