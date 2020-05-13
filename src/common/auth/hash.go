package auth

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/url"

	"github.com/wujiu2020/strip/utils"
)

func CreateHashValue(secret string, params map[string]string, now int64) (raw string) {
	h := md5.New()
	v := url.Values{}
	for key, value := range params {
		v.Set(key, value)
	}
	v.Set("t", utils.ValueTo(now).String())
	h.Write([]byte(v.Encode() + secret))
	hash := hex.EncodeToString(h.Sum(nil))[:10]
	v.Set("h", hash)
	raw = v.Encode()
	return
}

func ParseHashValue(secret, raw string, keys []string) (params map[string]string, now int64, err error) {
	v, err := url.ParseQuery(raw)
	if err != nil {
		return
	}
	params = make(map[string]string, len(keys))
	for _, key := range keys {
		params[key] = v.Get(key)
	}
	now = utils.StrTo(v.Get("t")).MustInt64()
	if CreateHashValue(secret, params, now) != v.Encode() {
		err = errors.New("wrong hash")
		return
	}
	return
}

func CreateHashValues(secret string, params url.Values, now int64) (raw string, err error) {
	var body string
	if len(secret) > 20 {
		secret = secret[:10]
		body = secret[10:]
	}

	h := hmac.New(md5.New, []byte(secret))
	_, err = h.Write([]byte(body))
	if err != nil {
		return
	}

	v := url.Values{}
	for key, value := range params {
		v[key] = value
	}

	v.Set("_t", utils.ValueTo(now).String())
	_, err = h.Write([]byte(v.Encode()))
	if err != nil {
		return
	}

	hash := hex.EncodeToString(h.Sum(nil))
	v.Set("_h", hash)
	raw = v.Encode()
	return
}

func ParseHashValues(secret, raw string) (params url.Values, now int64, err error) {
	params, err = url.ParseQuery(raw)
	if err != nil {
		return
	}
	now = utils.StrTo(params.Get("_t")).MustInt64()
	params.Del("_t")
	params.Del("_h")

	res, err := CreateHashValues(secret, params, now)
	if err != nil {
		return
	}
	if res != raw {
		err = errors.New("wrong hash")
		return
	}
	return
}

func ParseHashValuesExpired(secret, raw string, now int64) (params url.Values, duration int, err error) {
	params, before, err := ParseHashValues(secret, raw)
	if err != nil {
		return
	}

	duration = int(now - before)
	if duration < 0 {
		duration = 0
	}
	return
}

func CreateExpiredHash(secret string, params url.Values, now int64) (raw string, err error) {
	v := url.Values{}
	for key, value := range params {
		v[key] = value
	}
	v.Set("_t", utils.ValueTo(now).String())

	h := hmac.New(md5.New, []byte(secret))
	_, err = h.Write([]byte(v.Encode()))
	if err != nil {
		return
	}

	ret := url.Values{}
	ret.Set("h", hex.EncodeToString(h.Sum(nil)))
	ret.Set("t", utils.ValueTo(now).String())

	raw = ret.Encode()
	return
}

func VerifyExpiredHash(secret string, params url.Values, raw string, now int64) (duration int, err error) {
	ret, err := url.ParseQuery(raw)
	if err != nil {
		return
	}
	before := utils.StrTo(ret.Get("t")).MustInt64()

	res, err := CreateExpiredHash(secret, params, before)
	if err != nil {
		return
	}

	if res != raw {
		err = errors.New("wrong hash")
		return
	}
	duration = int(now - before)
	if duration < 0 {
		duration = 0
	}
	return
}
