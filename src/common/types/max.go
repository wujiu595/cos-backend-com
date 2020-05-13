package types

import (
	"time"
)

// MaxString return max value
func MaxString(vs ...string) string {
	if len(vs) == 0 {
		panic("string values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxInt return max value
func MaxInt(vs ...int) int {
	if len(vs) == 0 {
		panic("int values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxUInt return max value
func MaxUInt(vs ...uint) uint {
	if len(vs) == 0 {
		panic("uint values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxInt64 return max value
func MaxInt64(vs ...int64) int64 {
	if len(vs) == 0 {
		panic("int64 values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxUInt64 return max value
func MaxUInt64(vs ...uint64) uint64 {
	if len(vs) == 0 {
		panic("uint64 values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxInt32 return max value
func MaxInt32(vs ...int32) int32 {
	if len(vs) == 0 {
		panic("int32 values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxUInt32 return max value
func MaxUInt32(vs ...uint32) uint32 {
	if len(vs) == 0 {
		panic("uint32 values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxInt16 return max value
func MaxInt16(vs ...int16) int16 {
	if len(vs) == 0 {
		panic("int16 values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxUInt16 return max value
func MaxUInt16(vs ...uint16) uint16 {
	if len(vs) == 0 {
		panic("uint16 values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxInt8 return max value
func MaxInt8(vs ...int8) int8 {
	if len(vs) == 0 {
		panic("int8 values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxUInt8 return max value
func MaxUInt8(vs ...uint8) uint8 {
	if len(vs) == 0 {
		panic("uint8 values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxFloat32 return max value
func MaxFloat32(vs ...float32) float32 {
	if len(vs) == 0 {
		panic("float32 values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxFloat64 return max value
func MaxFloat64(vs ...float64) float64 {
	if len(vs) == 0 {
		panic("float64 values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxFloat64P3 return max value
func MaxFloat64P3(vs ...Float64P3) Float64P3 {
	if len(vs) == 0 {
		panic("Float64P3 values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

// MaxTime return max value
func MaxTime(vs ...time.Time) time.Time {
	if len(vs) == 0 {
		panic("time values required")
	}

	max := vs[0]
	for _, v := range vs[1:] {
		if v.After(max) {
			max = v
		}
	}

	return max
}
