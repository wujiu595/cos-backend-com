package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONArray is a json.RawMessage, which is a []byte underneath.
// Value() validates the json format in the source, and returns an error if
// the json is not valid.  Scan does no validation.  JSONArray additionally
// implements `Unmarshal`, which unmarshals the json within to an interface{}
type JSONArray json.RawMessage

var emptyJSONArray = JSONArray("[]")

// MarshalJSON returns the *j as the JSON encoding of j.
func (j JSONArray) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return emptyJSONArray, nil
	}
	return j, nil
}

// UnmarshalJSON sets *j to a copy of data
func (j *JSONArray) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("JSONArray: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// Value returns j as a value.  This does a validating unmarshal into another
// RawMessage.  If j is invalid json, it returns an error.
func (j JSONArray) Value() (driver.Value, error) {
	var m json.RawMessage
	var err = j.Unmarshal(&m)
	if err != nil {
		return []byte{}, err
	}
	return []byte(j), nil
}

// Scan stores the src in *j.  No validation is done.
func (j *JSONArray) Scan(src interface{}) error {
	var source []byte
	switch t := src.(type) {
	case string:
		source = []byte(t)
	case []byte:
		if len(t) == 0 {
			source = emptyJSONArray
		} else {
			source = t
		}
	case nil:
		*j = emptyJSONArray
	default:
		return errors.New("Incompatible type for JSONArray")
	}
	*j = JSONArray(append((*j)[0:0], source...))
	return nil
}

// Unmarshal unmarshal's the json in j to v, as in json.Unmarshal.
func (j *JSONArray) Unmarshal(v interface{}) error {
	if len(*j) == 0 {
		*j = emptyJSONArray
	}
	return json.Unmarshal([]byte(*j), v)
}

// String supports pretty printing for JSONArray types.
func (j JSONArray) String() string {
	if len(j) == 0 {
		return string(emptyJSONArray)
	}
	return string(j)
}
