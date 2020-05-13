package util

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func HeaderSetContentDisposition(header http.Header, filename string) {
	filename = strings.Replace(url.QueryEscape(filename), "+", "%20", -1)
	header.Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"; filename*=utf-8''%s`, filename, filename),
	)
}
