package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONMapString map[string]string

var emptyJSONMapString = []byte(`null`)
var _JSONMapString = JSONMapString{}
var _ json.Marshaler = _JSONMapString
var _ json.Unmarshaler = &_JSONMapString

// MarshalJSON implements the json.Marshaler interface.
func (j JSONMapString) MarshalJSON() ([]byte, error) {
	if j == nil {
		return emptyJSONMapString, nil
	}
	return json.Marshal(map[string]string(j))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (j *JSONMapString) UnmarshalJSON(b []byte) error {
	if j == nil {
		return errors.New("JSONMapString: UnmarshalJSON on nil pointer")
	}
	if b == nil {
		return nil
	}
	v := map[string]string{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	*j = v
	return nil
}

// Value implements the driver.Valuer interface.
func (j JSONMapString) Value() (driver.Value, error) {
	b, err := j.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Scan implements the sql.Scanner interface.
func (j *JSONMapString) Scan(src interface{}) error {
	var source []byte
	switch t := src.(type) {
	case string:
		source = []byte(t)
	case []byte:
		source = t
	case nil:
		return nil
	default:
		return errors.New("Incompatible type for JSONMapString")
	}
	if len(source) == 0 {
		return nil
	}
	return j.UnmarshalJSON(source)
}

func (j JSONMapString) String() string {
	b, err := j.MarshalJSON()
	if err != nil {
		return err.Error()
	}
	return string(b)
}
