package providers

import (
	"net/http"

	"github.com/wujiu2020/strip"
)

func LoadClassic(stp *strip.Strip) {
	rt := ClientRoundTripper()
	stp.Provide(rt)
	stp.ProvideAs(rt, (*http.RoundTripper)(nil))
	stp.Provide(HttpClient())
}
