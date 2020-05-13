package util

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	json "cos-backend-com/src/common/pgencoding/json2"
)

type PgJsonScanWrapValues []interface{}

func (p PgJsonScanWrapValues) Scan(src interface{}) (err error) {
	for i, value := range p {
		wrap := &PgJsonScanWrap{value}
		err = wrap.Scan(src)
		if err != nil {
			err = fmt.Errorf("scan index %d, error %v", i, err)
			break
		}
	}
	return
}

// PgJsonScanWrap 中的 struct tag 使用的是 `db:"name"`
// time should be type: TIMESTAMP WITH TIME ZONE
type PgJsonScanWrap struct {
	Value interface{}
}

func (p PgJsonScanWrap) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	data, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("JsonScanWrap must be a []byte, got %T instead", src)
	}
	var err error
	if value, ok := p.Value.(sql.Scanner); ok {
		err = value.Scan(src)
	} else {
		err = json.Unmarshal(data, p.Value)
	}
	return err
}

func PgMapQuery(sql string, maps map[string]interface{}) (string, []interface{}) {
	return pgMapQuery(sql, maps, true)
}

func PgMapQueryV2(sql string, maps map[string]interface{}) (string, []interface{}) {
	return pgMapQuery(sql, maps, false)
}

func pgMapQuery(sql string, maps map[string]interface{}, includeArrayBrace bool) (string, []interface{}) {
	var (
		args = make([]interface{}, 0, len(maps))
		n    = 0
	)
	for name, val := range maps {
		placeholder := "$" + name
		if strings.Contains(sql, placeholder) {
			rval := reflect.ValueOf(val)
			_, ok := val.(driver.Valuer)
			if !ok && rval.Kind() == reflect.Slice {
				num := rval.Len() - 1
				params := ""
				for i := 0; i <= num; i++ {
					n++
					if i < num {
						params += fmt.Sprintf("$%d, ", n)
					}
					args = append(args, rval.Index(i).Interface())
				}
				if num >= 0 {
					params += fmt.Sprintf("$%d", n)
				}
				if includeArrayBrace {
					sql = strings.Replace(sql, placeholder, "("+params+")", -1)
				} else {
					sql = strings.Replace(sql, placeholder, params, -1)
				}
				continue
			}
			n++
			sql = strings.Replace(sql, placeholder, fmt.Sprintf("$%d", n), -1)
			args = append(args, val)
		}
	}
	return sql, args
}

func PgJsonScanWrapGetSlices(val reflect.Value, vals *[]reflect.Value, n int) {
	itf := val.Interface()
	if v, ok := itf.(PgJsonScanWrap); ok {
		PgJsonScanWrapGetSlices(reflect.Indirect(reflect.ValueOf(v.Value)), vals, n)
		return
	}
	if vs, ok := itf.(PgJsonScanWrapValues); ok {
		for _, v := range vs {
			PgJsonScanWrapGetSlices(reflect.Indirect(reflect.ValueOf(v)), vals, n)
		}
		return
	}
	if val.Kind() == reflect.Slice && val.Len() == n {
		*vals = append(*vals, val)
	}
}

var likeReplacer = strings.NewReplacer(
	`\`, `\\`,
	`%`, `\%`,
	`_`, `\_`,
)

func PgEscapeLike(value string) string {
	return likeReplacer.Replace(value)
}
