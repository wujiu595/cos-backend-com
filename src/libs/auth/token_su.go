package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"net/url"

	"github.com/wujiu2020/strip/utils"
)

func CreateSuInfo(suid int64, accessToken string) (suInfo string, err error) {
	values := url.Values{}
	values.Set("suid", utils.ToStr(suid))
	h := hmac.New(sha1.New, []byte(accessToken))
	_, err = h.Write([]byte(values.Encode()))
	if err != nil {
		return
	}
	hash := hex.EncodeToString(h.Sum(nil))
	values.Set("h", hash)
	suInfo = base64.StdEncoding.EncodeToString([]byte(values.Encode()))
	return
}
