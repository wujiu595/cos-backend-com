package types

import (
	"database/sql/driver"
	"encoding/json"
)

type JSONAny struct {
	Any interface{}
}

// MarshalJSON returns the *j as the JSON encoding of j.
func (j JSONAny) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Any)
}

func (j JSONAny) Value() (driver.Value, error) {
	return j.MarshalJSON()
}

func (j JSONAny) String() string {
	b, err := j.MarshalJSON()
	if err != nil {
		return err.Error()
	}
	return string(b)
}
