package users

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/auth"
	"cos-backend-com/src/libs/sdk/account"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/wujiu2020/strip/sessions"
	"github.com/wujiu2020/strip/utils"
)

func createOAuthToken(id, ak, sk string) (token *auth.OAuth2Token, err error) {
	values := url.Values{}
	values.Set("id", id)
	values.Set("exp", utils.ToStr(refreshExpiresIn))
	refreshToken, err := createToken(values.Encode(), TokenTypeRefreshToken, ak, sk)
	if err != nil {
		return
	}

	accessToken, err := createAccessToken(id, ak, sk)
	if err != nil {
		return
	}

	token = &auth.OAuth2Token{}
	token.AccessToken = accessToken
	token.RefreshToken = refreshToken
	token.ExpiresIn = tokenExpiresIn
	return token, nil
}

func createAccessToken(id, ak, sk string) (accessToken string, err error) {
	values := url.Values{}
	values.Set("id", id)
	values.Set("exp", utils.ToStr(tokenExpiresIn))

	accessToken, err = createToken(values.Encode(), TokenTypeAccessToken, ak, sk)
	return
}

func createToken(params, salt, ak, sk string) (string, error) {
	value, ok := sessions.EncodeSecureValue(params, ak+sk, time.Now())
	if !ok {
		return "", apierror.ErrApiAuthCreateFailed
	}

	token := value
	h := hmac.New(sha1.New, []byte(salt))
	h.Write([]byte(token))
	token += ":" + hex.EncodeToString(h.Sum(nil))
	token = base64.StdEncoding.EncodeToString([]byte(token))
	return token, nil
}

func (p *accessTokens) nextId(ctx context.Context) (netxtId flake.ID, err error) {
	err = p.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &netxtId, `SELECT nextval('global_id_sequence')`)
	})
	return
}

func verifyTokenOutter(token string, salt account.TokenType) (raw string, err error) {
	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return
	}
	token = string(b)

	parts := strings.SplitN(token, ":", 2)
	if len(parts) != 2 {
		err = errors.New("token verify: invalid size")
		return
	}

	h := hmac.New(sha1.New, []byte(salt))
	h.Write([]byte(parts[0]))
	if hex.EncodeToString(h.Sum(nil)) != parts[1] {
		err = errors.New("token verify: invalid token")
		return
	}

	return parts[0], nil
}

func verifyTokenInner(token string, exp int64, ak, sk string) (err error) {
	_, createAt, ok := sessions.DecodeSecureValue(token, ak+sk)
	if !ok {
		err = apierror.ErrApiAuthInvalidToken
		return
	}

	if createAt.Add(time.Duration(exp) * time.Second).Before(time.Now()) {
		err = apierror.ErrApiAuthExpiredToken
		return
	}
	return
}

func getValuesFromToken(token string) (id string, exp int64, ok bool) {
	raw := sessions.GetRawFromSecureValue(token)
	values, err := url.ParseQuery(raw)
	if err != nil {
		return
	}
	id = values.Get("id")
	exp, _ = utils.StrTo(values.Get("exp")).Int64()
	ok = true
	return
}
