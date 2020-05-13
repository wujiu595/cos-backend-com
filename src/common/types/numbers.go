package types

import "strconv"

type Float64P3 float64

func (p Float64P3) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(p), 'f', 3, 64)), nil
}
