package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONMapAny map[string]interface{}

var emptyJSONMapAny = []byte(`{}`)
var _JSONMapAny = JSONMapAny{}
var _ json.Marshaler = _JSONMapAny
var _ json.Unmarshaler = &_JSONMapAny

// MarshalJSON implements the json.Marshaler interface.
func (j JSONMapAny) MarshalJSON() ([]byte, error) {
	if j == nil {
		return emptyJSONMapAny, nil
	}
	return json.Marshal(map[string]interface{}(j))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (j *JSONMapAny) UnmarshalJSON(b []byte) error {
	if j == nil {
		return errors.New("JSONMapAny: UnmarshalJSON on nil pointer")
	}
	if b == nil {
		return nil
	}
	v := map[string]interface{}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	*j = v
	return nil
}

// Value implements the driver.Valuer interface.
func (j JSONMapAny) Value() (driver.Value, error) {
	b, err := j.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Scan implements the sql.Scanner interface.
func (j *JSONMapAny) Scan(src interface{}) error {
	var source []byte
	switch t := src.(type) {
	case string:
		source = []byte(t)
	case []byte:
		source = t
	case nil:
		return nil
	default:
		return errors.New("Incompatible type for JSONMapAny")
	}
	if len(source) == 0 {
		return nil
	}
	return j.UnmarshalJSON(source)
}

func (j JSONMapAny) String() string {
	b, err := j.MarshalJSON()
	if err != nil {
		return err.Error()
	}
	return string(b)
}
