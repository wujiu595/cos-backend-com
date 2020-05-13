package auth

import (
	"fmt"
	"net/url"
	"time"

	"github.com/wujiu2020/strip/caches"
	"github.com/wujiu2020/strip/utils"

	"cos-backend-com/src/common/flake"
)

func CreateSuToken(cache caches.CacheProvider, secret string, uid flake.ID) (token string, err error) {
	key := fmt.Sprintf("%s:%v:%s", "su", uid, utils.MustRandomCreateString(10))
	values := url.Values{"id": {key}}
	err = cache.Set(key, uid.String(), 60)
	if err != nil {
		return
	}
	token = CreateParamsToken(secret, values)
	return
}

func ParseSuTokenOnce(cache caches.CacheProvider, secret string, token string) (uid flake.ID, err error) {
	params, err := ParseParamsToken(secret, token, time.Now().Unix(), 60)
	if err != nil {
		return
	}
	key := params.Get("id")
	defer cache.Delete(key)
	if v, er := cache.Get(key); er != nil {
		err = er
		return
	} else {
		uid = flake.MustFromString(v.String())
	}
	return
}
