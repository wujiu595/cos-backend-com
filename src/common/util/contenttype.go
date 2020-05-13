package util

import (
	"net/http"
	"strings"
)

const (
	contentTypeJson = "application/json"
)

func IsReqAcceptJson(req *http.Request) bool {
	accept := req.Header.Get("Accept")
	return strings.Contains(accept, contentTypeJson)
}
