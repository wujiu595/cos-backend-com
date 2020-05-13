package schema

import (
	"net/url"

	"github.com/go-playground/form"
)

var decoder = form.NewDecoder()
var encoder = form.NewEncoder()

func Decode(v interface{}, values url.Values) (err error) {
	err = decoder.Decode(v, values)
	return
}

func Encode(v interface{}) (values url.Values, err error) {
	values, err = encoder.Encode(v)
	return
}
