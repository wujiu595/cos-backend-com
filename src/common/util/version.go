package util

import (
	"net/http"

	"github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/utils/helpers"
)

var (
	versionDefault = "null"
	Version        = "null"
)

func HasVersion() bool {
	return Version != versionDefault && Version != ""
}

func VersionHandler(req *http.Request, resp http.ResponseWriter) {
	resp.Write([]byte(Version))
}

func PrintVersion() {
	helpers.X.Infof("version: %v", Version)
}

func VersionRouter() strip.Handler {
	PrintVersion()
	return strip.Router("/strip", strip.Get(VersionHandler))
}
