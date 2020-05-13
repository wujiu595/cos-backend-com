package dbcache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	dbjson "cos-backend-com/src/common/pgencoding/json2"
)

var (
	DefaultTimeout = time.Minute

	ErrDataNotFound = fmt.Errorf("data not found")
	ErrDecodeFailed = fmt.Errorf("decode to struct failed")
)

func (p *Table) Datas(idValues ...interface{}) []*Data {
	datas := make([]*Data, 0, len(idValues))
	for _, v := range idValues {
		datas = append(datas, p.Data(v))
	}
	return datas
}

func (p *Table) Data(idValue interface{}) *Data {
	return &Data{Table: p, IDValue: fmt.Sprint(idValue)}
}

func (p *Table) buildQuery(lth int) string {
	selectSQL := "*"
	if len(p.cfg.Cols) > 0 {
		selectSQL = `t.` + strings.Join(p.cfg.Cols, `, t.`)
	}

	whereSQL := ""
	if lth > 0 {
		whereSQL += ` AND t.` + p.cfg.IdName + ` IN (`
		whereSQL += "$1"
		for i := 2; i <= lth; i++ {
			whereSQL += ", $" + fmt.Sprint(i)
		}
		whereSQL += `)`
	}
	return `SELECT ` + selectSQL + ` FROM ` + p.cfg.Table + ` AS t WHERE 1=1` + whereSQL
}

func (p *Data) key() string {
	return p.Table.prefix + ":" + p.IDValue
}

func getTimeoutDur(params ...int) time.Duration {
	var timeout time.Duration
	if len(params) > 0 {
		switch {
		case params[0] > 0:
			timeout = time.Duration(params[0]) * time.Second
			return timeout
		default:
			timeout = time.Duration(params[0])
		}
	}

	if timeout == 0 {
		timeout = DefaultTimeout
	}
	if timeout < 0 {
		timeout = 0
	}
	return timeout
}

func sqlxMapMarshalToStruct(m, v interface{}) (b []byte, err error) {
	isBin := false
	switch v := m.(type) {
	case []byte:
		isBin = true
		b = v
	case map[string]interface{}:
		for key, value := range v {
			if d, ok := value.([]byte); ok {
				v[key] = json.RawMessage(d)
			}
		}
	default:
		panic(fmt.Errorf("unsupport m type `%T` in sqlxMapMarshalToStruct", m))
	}
	if !isBin {
		b, err = json.Marshal(m)
		if err != nil {
			return
		}
	}
	err = dbjson.Unmarshal(b, v)
	return
}

var testPrefixes = map[string]bool{}

func encodePrefix(name string) string {
	h := md5.New()
	h.Write([]byte(name))
	v := hex.EncodeToString(h.Sum(nil))[:4]
	if _, ok := testPrefixes[v]; ok {
		panic("not support repeat defined table")
	}
	testPrefixes[v] = true
	return v
}

type Fetcher func(key1, key2, output interface{}) error

func GetValue(redis Header, prefix string, key1, key2, out interface{}, fetcher Fetcher) (err error) {
	key := fmt.Sprintf("%v_%v_%v", prefix, key1, key2)
	res, err := redis.MultiGet([]string{key})
	if err != nil {
		return
	}

	missed := false
	for _, b := range res {
		if b != nil {
			er := dbjson.Unmarshal(b, out)
			if er != nil {
				return
			} else {
				break
			}
		}
		missed = true
	}
	if !missed {
		return
	}

	err = fetcher(key1, key2, out)
	if err != nil {
		return
	}

	keys := make([]string, 0, 1)
	bs := make([][]byte, 0, 1)
	expires := make([]int, 0, 1)
	keys = append(keys, key)
	data, err := dbjson.Marshal(out)
	if err != nil {
		return
	}
	bs = append(bs, data)
	expires = append(expires, 300)
	err = redis.MultiSet(keys, bs, expires)
	return
}
