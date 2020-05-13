package routers

import (
	"net/http"

	"github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/params"
)

type Base struct {
	Req *http.Request       `inject`
	Rw  http.ResponseWriter `inject`

	Params *params.Params `inject`
	Log    strip.Logger   `inject`

	Ctx strip.Context `inject`
}
