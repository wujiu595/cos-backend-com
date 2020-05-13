package filters

import (
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/utils"
)

func OpenTracingFilter() interface{} {
	opNameFunc := func(req *http.Request) string {
		return "HTTP " + req.Method
	}
	spanObserver := func(span opentracing.Span, r *http.Request) {}
	return func(ctx strip.Context, log strip.ReqLogger, tr opentracing.Tracer, rw http.ResponseWriter, req *http.Request) {

		spanCtx, _ := tr.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		sp := tr.StartSpan(opNameFunc(req), ext.RPCServerOption(spanCtx))
		ext.HTTPMethod.Set(sp, req.Method)
		ext.HTTPUrl.Set(sp, req.URL.String())
		spanObserver(sp, req)

		reqCtx := utils.CtxWithValue(req.Context(), &tr)
		req = req.WithContext(reqCtx)
		req = req.WithContext(opentracing.ContextWithSpan(req.Context(), sp))

		ctx.Provide(req)
		ctx.ReplaceContext(req.Context())

		var status uint16 = 200

		if rw, ok := rw.(strip.ResponseWriter); ok {
			rw.Before(func(strip.ResponseWriter) {
				if rw.Status() != 0 {
					status = uint16(rw.Status())
				}

				ext.HTTPStatusCode.Set(sp, status)

				// TODO ?? 需要增加其他 request 信息
				sp.SetTag("http.x-reqid", log.ReqId())

				sp.Finish()
			})
		}
	}
}
