package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

var timeLayout = "15:04:05"

type TimeShift struct {
	time.Time
}

func (ts TimeShift) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ts.Format(timeLayout) + `"`), nil
}

func (ts *TimeShift) UnmarshalJSON(b []byte) error {
	s := string(b)
	// len(`"23:59:00"`) == 10
	if len(s) != 10 {
		return fmt.Errorf(`TimeParseError: should be a string formatted as "15:04:05, %s`, s)
	}
	ret, err := time.Parse(timeLayout, s[1:9])
	if err != nil {
		return err
	}
	ts.Time = ret
	return nil
}

func (ts TimeShift) Value() (driver.Value, error) {
	return ts.Time.Format(timeLayout), nil
}

func (ts *TimeShift) Scan(src interface{}) error {
	skip := false
	var source []byte
	switch t := src.(type) {
	case string:
		source = []byte(t)
	case []byte:
		source = t
	case time.Time:
		ts.Time = t
		skip = true
	default:
		return errors.New("Incompatible type for TimeShift")
	}

	if !skip {
		return ts.UnmarshalJSON(source)
	}

	return nil
}

// String supports pretty printing for JSONText types.
func (ts TimeShift) String() string {
	return ts.Time.Format(timeLayout)
}
