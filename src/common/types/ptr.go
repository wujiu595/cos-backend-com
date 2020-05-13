package types

import (
	"time"

	"cos-backend-com/src/common/flake"
)

// StringP return pointer
func StringP(v string) *string { return &v }

// IntP return pointer
func IntP(v int) *int { return &v }

// UIntP return pointer
func UIntP(v uint) *uint { return &v }

// Int64P return pointer
func Int64P(v int64) *int64 { return &v }

// UInt64P return pointer
func UInt64P(v uint64) *uint64 { return &v }

// Int32P return pointer
func Int32P(v int32) *int32 { return &v }

// UInt32P return pointer
func UInt32P(v uint32) *uint32 { return &v }

// Int16P return pointer
func Int16P(v int16) *int16 { return &v }

// UInt16P return pointer
func UInt16P(v uint16) *uint16 { return &v }

// Int8P return pointer
func Int8P(v int8) *int8 { return &v }

// UInt8P return pointer
func UInt8P(v uint8) *uint8 { return &v }

// Float32P return pointer
func Float32P(v float32) *float32 { return &v }

// Float64P return pointer
func Float64P(v float64) *float64 { return &v }

// Float64P3P return pointer
func Float64P3P(v Float64P3) *Float64P3 { return &v }

// TimeP return pointer
func TimeP(v time.Time) *time.Time { return &v }

// BoolP return pointer
func BoolP(v bool) *bool { return &v }

// FlakeP returns pointer
func FlakeP(v flake.ID) *flake.ID { return &v }
