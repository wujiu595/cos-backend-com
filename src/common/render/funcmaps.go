package render

import (
	"encoding/json"
	"fmt"
	"text/template"
)

var defaultFuncMaps = template.FuncMap{
	"json_string": func(value interface{}) string {
		b, _ := json.Marshal(fmt.Sprint(value))
		if len(b) >= 2 {
			return string(b[1 : len(b)-1])
		}
		return ""
	},
}
