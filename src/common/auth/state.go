package auth

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"time"
)

var (
	ErrStateExpired = fmt.Errorf("token expired")
)

func CreateStateToken(secret string, params url.Values) string {
	return CreateStateTokenN(secret, params, time.Now().Unix())
}

func VerifyStateToken(secret string, params url.Values, raw string, limitRemain int) error {
	return VerifyStateTokenN(secret, params, raw, time.Now().Unix(), limitRemain)
}

func CreateStateTokenN(secret string, params url.Values, now int64) string {
	raw, err := CreateExpiredHash(secret, params, now)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func VerifyStateTokenN(secret string, params url.Values, raw string, now int64, limitDuration int) error {
	b, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return err
	}
	raw = string(b)
	duration, err := VerifyExpiredHash(secret, params, raw, now)
	if err != nil {
		return err
	}
	if duration > limitDuration {
		return ErrStateExpired
	}
	return nil
}

func CreateParamsToken(secret string, params url.Values) string {
	raw, err := CreateHashValues(secret, params, time.Now().Unix())
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func ParseParamsToken(secret, raw string, now int64, limitDuration int) (params url.Values, err error) {
	b, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return
	}
	raw = string(b)
	params, duration, err := ParseHashValuesExpired(secret, raw, now)
	if err != nil {
		return
	}
	if duration > limitDuration {
		err = ErrStateExpired
		return
	}
	return
}
