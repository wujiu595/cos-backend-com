package files

import (
	"net/http"
)

func NewFileService(conf BaseConfig) interface{} {
	return func(rt http.RoundTripper) FileService {
		return NewClient(ClientConfig{
			conf,
			rt,
		})
	}
}
