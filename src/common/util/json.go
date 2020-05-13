package util

import (
	"encoding/json"
)

func JsonMapToStruct(inf interface{}, val interface{}) (err error) {
	b, err := json.Marshal(inf)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, val)
	return
}
