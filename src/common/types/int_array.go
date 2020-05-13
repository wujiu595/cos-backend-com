package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

// IntArray represents a one-dimensional array of the PostgreSQL integer types.
type IntArray []int

// Scan implements the sql.Scanner interface.
func (a *IntArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to IntArray", src)
}

func (a *IntArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "IntArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(IntArray, len(elems))
		for i, v := range elems {
			if b[i], err = strconv.Atoi(string(v)); err != nil {
				return fmt.Errorf("pq: parsing array element index %d: %v", i, err)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a IntArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendInt(b, int64(a[0]), 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendInt(b, int64(a[i]), 10)
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}
