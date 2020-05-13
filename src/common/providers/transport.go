package providers

import (
	"net/http"
	"time"

	"github.com/wujiu2020/strip"

	"cos-backend-com/src/common/client"
)

func ClientRoundTripper() interface{} {
	return func(log strip.Logger) client.RoundTripper {
		return client.NewTransportWithLogger(log).With(client.NewOpenTracingTransport(nil))
	}
}

func HttpClient() interface{} {
	return func(tr http.RoundTripper) *http.Client {
		return &http.Client{Transport: tr, Timeout: time.Minute * 2}
	}
}
