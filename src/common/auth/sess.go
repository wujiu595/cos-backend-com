package auth

import (
	"net/http"
)

func ReadParamsFromCookie(secret string, keys []string, cookieName string, req *http.Request) (params map[string]string, now int64, err error) {
	cookie, err := req.Cookie(cookieName)
	if err != nil {
		return
	}

	params, now, err = ParseHashValue(secret, cookie.Value, keys)
	return
}

func WriteParamsToCookie(secret string, params map[string]string, cookieName, cookieDomain string, now int64, rw http.ResponseWriter) {
	curUser := CreateHashValue(secret, params, now)
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  curUser,
		Path:   "/",
		Domain: cookieDomain,
		MaxAge: 365 * 24 * 3600,
	}
	http.SetCookie(rw, cookie)
}
