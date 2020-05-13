package kv

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"cos-backend-com/src/common/proto"
)

var (
	Codes = map[byte][]byte{
		'\t': []byte(`\\t`),
		'"':  []byte(`\"`),
		' ':  []byte(`\ `),
		'=':  []byte(`\=`),
	}

	codesStr = map[string]string{}
)

func Decode(bytes []byte) (data proto.Data, err error) {
	line := string(bytes)
	data = make(proto.Data)

	len := len(line)
	if len == 0 {
		return
	}

	arr := []byte(line)
	start, cur, del_pos := 0, 0, 0
	kv_del := byte('=')
	field_del := byte('\t')
	seenFiledDelimeter, seenKVDelimeter, seenChar := false, false, false

	for cur < len {
		switch arr[cur] {
		case kv_del:
			if !seenKVDelimeter {
				seenKVDelimeter = true
				del_pos = cur
			}
		case field_del:
			seenFiledDelimeter = true
		default:
			seenChar = true
		}

		if seenFiledDelimeter {
			if start < del_pos && del_pos <= cur {
				key := string(arr[start:del_pos])
				val := string(arr[del_pos+1 : cur])
				data[key] = val
			} else {
				if seenChar {
					err = fmt.Errorf("no delimeter found or empty key/value at %v %v %v", start, del_pos, cur)
					return
				}
			}

			start = cur + 1
			seenFiledDelimeter = false
			seenKVDelimeter = false
			seenChar = false
		}
		cur++
	}

	if start < del_pos && del_pos < cur {
		key := string(arr[start:del_pos])
		val := string(arr[del_pos+1 : len])
		data[key] = decodeStringField(val)
	}

	return
}

func Encode(data proto.Data) (bytes []byte) {
	b := []byte{}
	keys := make([]string, len(data))
	i := 0
	for k := range data {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := data[k]
		b = append(b, []byte(String(k))...)
		b = append(b, '=')
		switch t := v.(type) {
		case int:
			b = append(b, []byte(strconv.FormatInt(int64(t), 10))...)
		case int8:
			b = append(b, []byte(strconv.FormatInt(int64(t), 10))...)
		case int16:
			b = append(b, []byte(strconv.FormatInt(int64(t), 10))...)
		case int32:
			b = append(b, []byte(strconv.FormatInt(int64(t), 10))...)
		case int64:
			b = append(b, []byte(strconv.FormatInt(t, 10))...)
		case uint:
			b = append(b, []byte(strconv.FormatInt(int64(t), 10))...)
		case uint8:
			b = append(b, []byte(strconv.FormatInt(int64(t), 10))...)
		case uint16:
			b = append(b, []byte(strconv.FormatInt(int64(t), 10))...)
		case uint32:
			b = append(b, []byte(strconv.FormatInt(int64(t), 10))...)
		case uint64:
			val := []byte(strconv.FormatFloat(float64(t), 'f', -1, 64))
			b = append(b, val...)
		case float32:
			val := []byte(strconv.FormatFloat(float64(t), 'f', -1, 32))
			b = append(b, val...)
		case float64:
			val := []byte(strconv.FormatFloat(t, 'f', -1, 64))
			b = append(b, val...)
		case bool:
			b = append(b, []byte(strconv.FormatBool(t))...)
		case []byte:
			b = append(b, t...)
		case string:
			b = append(b, []byte(escapeStringField(t))...)
		case nil:
			// skip
		default:
			// Can't determine the type, so convert to string
			b = append(b, []byte(escapeStringField(fmt.Sprintf("%v", v)))...)

		}
		b = append(b, '\t')
	}
	if len(b) > 0 {
		return b[0 : len(b)-1]
	}
	return b
}

func String(in string) string {
	for b, esc := range codesStr {
		in = strings.Replace(in, b, esc, -1)
	}
	return in
}

func escapeStringField(in string) string {
	var out []byte
	i := 0
	for {
		if i >= len(in) {
			break
		}
		// escape double-quotes
		if in[i] == '\\' {
			out = append(out, '\\')
			out = append(out, '\\')
			i++
			continue
		}
		// escape double-quotes
		if in[i] == '"' {
			out = append(out, '\\')
			out = append(out, '"')
			i++
			continue
		}
		// escape double-quotes
		if in[i] == '\t' {
			out = append(out, '\\')
			out = append(out, 't')
			i++
			continue
		}
		// escape double-quotes
		if in[i] == '\n' {
			out = append(out, '\\')
			out = append(out, 'n')
			i++
			continue
		}
		out = append(out, in[i])
		i++

	}
	return string(out)
}
func decodeStringField(in string) string {
	var out []byte
	i := 0
	for {
		if i >= len(in) {
			break
		}
		// escape double-quotes
		if in[i] == '\\' && i+1 < len(in) && (in[i+1] == '\\' || in[i+1] == '"' || in[i+1] == '\t' || in[i+1] == '\n') {
			out = append(out, in[i+1])
			i += 2
			continue
		}
		out = append(out, in[i])
		i++

	}
	return string(out)
}
