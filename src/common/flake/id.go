package flake

import (
	"database/sql/driver"
	"encoding"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type ID int64

var _id = ID(0)
var _ json.Marshaler = _id
var _ json.Unmarshaler = &_id
var _ encoding.TextMarshaler = _id
var _ encoding.TextUnmarshaler = &_id
var _ encoding.BinaryMarshaler = _id
var _ encoding.BinaryUnmarshaler = &_id

// MarshalJSON implements the json.Marshaler interface.
func (p ID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + p.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (p *ID) UnmarshalJSON(b []byte) error {
	if len(b) >= 3 && b[0] == '"' && b[len(b)-1] == '"' {
		return p.UnmarshalText(b[1 : len(b)-1])
	}
	return p.UnmarshalText(b)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (p ID) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (p *ID) UnmarshalText(b []byte) error {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid ID %q, err: %v", string(b), err)
	}

	*p = ID(i)
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (p ID) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (p *ID) UnmarshalBinary(b []byte) error {
	if len(b) != 8 {
		return fmt.Errorf("must be exactly 8 bytes long, got %d bytes", len(b))
	}
	*p = ID(binary.BigEndian.Uint64(b))
	return nil
}

// Value implements the driver.Valuer interface.
func (p ID) Value() (driver.Value, error) {
	return p.Int64(), nil
}

// Scan implements the sql.Scanner interface.
func (p *ID) Scan(src interface{}) error {
	switch v := src.(type) {
	case int64:
		*p = ID(v)
		return nil
	case []byte:
		if len(v) == 8 {
			return p.UnmarshalBinary(v)
		}
		return p.UnmarshalText(v)
	case string:
		return p.UnmarshalText([]byte(v))
	}
	return fmt.Errorf("cannot convert %T to ID", src)
}

func (p ID) Time(timeShiftBits uint8, epoch int64) time.Time {
	return time.Unix(0, ((p.Int64()>>timeShiftBits)+epoch)*1e6)
}

func (p ID) Int() int {
	return int(p)
}

func (p ID) Int64() int64 {
	return int64(p)
}

func (p ID) Bytes() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(p))
	return b
}

func (p ID) String() string {
	return strconv.FormatInt(p.Int64(), 10)
}

type IDS []ID

func (s IDS) Len() int {
	return len(s)
}
func (s IDS) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s IDS) Less(i, j int) bool {
	return s[i] < s[j]
}

func FromString(str string) (ID, error) {
	var id ID
	err := id.UnmarshalText([]byte(str))
	return id, err
}

func MustFromString(str string) ID {
	var id ID
	id.UnmarshalText([]byte(str))
	return id
}
