package pgbuilder

import (
	"reflect"
	"time"

	"cos-backend-com/src/common/types"

	"github.com/gocraft/dbr"
)

// Builder alias for dbr.Builder
type Builder = dbr.Builder

// conditions
var (
	Gt  = dbr.Gt
	Gte = dbr.Gte
	Lt  = dbr.Lt
	Lte = dbr.Lte
)

// operators
var (
	And = dbr.And
	Or  = dbr.Or
)

// Eq return equal express, when value is nil, build `IS NULL`
func Eq(column string, value interface{}) dbr.Builder {
	if IsNil(value) {
		return dbr.Eq(column, nil)
	}

	return dbr.Eq(column, value)
}

// Neq return equal express, when value is nil, build `IS NOT NULL`
func Neq(column string, value interface{}) dbr.Builder {
	if IsNil(value) {
		return dbr.Neq(column, nil)
	}

	return dbr.Neq(column, value)
}

// IsNil if the value is nil
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}

	switch vt := v.(type) {
	case string, int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint, float64, float32, types.Float64P3, time.Time:
		return false

	case *string:
		return vt == nil

	case *int8:
		return vt == nil

	case *uint8:
		return vt == nil

	case *int16:
		return vt == nil

	case *uint16:
		return vt == nil

	case *int32:
		return vt == nil

	case *uint32:
		return vt == nil

	case *int64:
		return vt == nil

	case *uint64:
		return vt == nil

	case *int:
		return vt == nil

	case *uint:
		return vt == nil

	case *float64:
		return vt == nil

	case *float32:
		return vt == nil

	case *types.Float64P3:
		return vt == nil

	case *time.Time:
		return vt == nil

	}

	val := reflect.ValueOf(v)
	k := val.Kind()
	return (k == reflect.Ptr || k == reflect.Slice) && val.IsNil()
}
