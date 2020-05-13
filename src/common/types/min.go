package types

import (
	"time"
)

// min

// MinString return min value
func MinString(vs ...string) string {
	if len(vs) == 0 {
		panic("string values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinInt return min value
func MinInt(vs ...int) int {
	if len(vs) == 0 {
		panic("int values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinUInt return min value
func MinUInt(vs ...uint) uint {
	if len(vs) == 0 {
		panic("uint values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinInt64 return min value
func MinInt64(vs ...int64) int64 {
	if len(vs) == 0 {
		panic("int64 values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinUInt64 return min value
func MinUInt64(vs ...uint64) uint64 {
	if len(vs) == 0 {
		panic("uint64 values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinInt32 return min value
func MinInt32(vs ...int32) int32 {
	if len(vs) == 0 {
		panic("int32 values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinUInt32 return min value
func MinUInt32(vs ...uint32) uint32 {
	if len(vs) == 0 {
		panic("uint32 values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinInt16 return min value
func MinInt16(vs ...int16) int16 {
	if len(vs) == 0 {
		panic("int16 values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinUInt16 return min value
func MinUInt16(vs ...uint16) uint16 {
	if len(vs) == 0 {
		panic("uint16 values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinInt8 return min value
func MinInt8(vs ...int8) int8 {
	if len(vs) == 0 {
		panic("int8 values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinUInt8 return min value
func MinUInt8(vs ...uint8) uint8 {
	if len(vs) == 0 {
		panic("uint8 values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinFloat32 return min value
func MinFloat32(vs ...float32) float32 {
	if len(vs) == 0 {
		panic("float32 values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinFloat64 return min value
func MinFloat64(vs ...float64) float64 {
	if len(vs) == 0 {
		panic("float64 values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinFloat64P3 return min value
func MinFloat64P3(vs ...Float64P3) Float64P3 {
	if len(vs) == 0 {
		panic("Float64P3 values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v < min {
			min = v
		}
	}

	return min
}

// MinTime return min value
func MinTime(vs ...time.Time) time.Time {
	if len(vs) == 0 {
		panic("time values required")
	}

	min := vs[0]
	for _, v := range vs[1:] {
		if v.Before(min) {
			min = v
		}
	}

	return min
}
