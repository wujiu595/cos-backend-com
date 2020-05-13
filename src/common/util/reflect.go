package util

import (
	"reflect"
)

func FindStructElemByName(elm reflect.Value, kind reflect.Kind, name string) (res reflect.Value, ok bool) {
	elm = reflect.Indirect(elm)
	val := reflect.Indirect(elm.FieldByName(name))
	if val.Kind() == kind {
		res = val
		ok = true
		return
	}
	return
}

func FindStructElemRecursive(elm reflect.Value, typ reflect.Type) (res reflect.Value, ok bool) {
	elm = reflect.Indirect(elm)
	if typ.Kind() != reflect.Struct {
		panic("typ must be a struct type")
	}
	if elm.Type() == typ {
		ok = true
		res = elm
		return
	}
	for i := 0; i < elm.NumField(); i++ {
		field := elm.Field(i)
		ftyp := field.Type()
		if ftyp.Kind() == reflect.Ptr {
			ftyp = ftyp.Elem()
		}
		if ftyp.Kind() != reflect.Struct {
			continue
		}

		if v, o := FindStructElemRecursive(field, typ); o {
			res = v
			ok = o
			return
		}
	}
	return
}

func CollectStructTags(elm reflect.Value, key string) []string {
	elm = reflect.Indirect(elm)
	typ := elm.Type()
	if typ.Kind() != reflect.Struct {
		panic("typ must be a struct type")
	}
	keys := make([]string, 0, elm.NumField())
	for i := 0; i < elm.NumField(); i++ {
		if v := typ.Field(i).Tag.Get(key); v != "" {
			keys = append(keys, v)
		}
	}
	return keys
}
