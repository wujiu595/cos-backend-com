package auth

import (
	"cos-backend-com/src/common/flake"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
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

func VerifySuInfo(raw string, accessToken string) (suid flake.ID, err error) {
	b, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return
	}

	values, err := url.ParseQuery(string(b))
	if err != nil {
		return
	}

	hash := values.Get("h")
	values.Del("h")

	h := hmac.New(sha1.New, []byte(accessToken))
	_, err = h.Write([]byte(values.Encode()))
	if err != nil {
		return
	}

	if hex.EncodeToString(h.Sum(nil)) != hash {
		err = fmt.Errorf("wrong su info hash")
		return
	}

	v, err := utils.StrTo(values.Get("suid")).Int64()
	if err != nil {
		return
	}
	suid = flake.ID(v)

	if suid <= 0 {
		err = fmt.Errorf("wrong suid in su info: %d", suid)
	}
	return
}
