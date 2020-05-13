package testutil

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// MarshalJson 出错即退出
func (p *IntegrationTest) MarshalJson(v interface{}) []byte {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		p.Log.Fatalf("MarshalJson: %s", err)
	}
	return b
}

func copyHeaders(header http.Header) http.Header {
	h := make(http.Header)
	for k, v := range header {
		p := make([]string, len(v))
		for i, a := range v {
			p[i] = a
		}
		h[k] = p
	}
	return h
}

func copyValues(values url.Values) url.Values {
	h := make(url.Values)
	for k, v := range values {
		p := make([]string, len(v))
		for i, a := range v {
			p[i] = a
		}
		h[k] = p
	}
	return h
}
