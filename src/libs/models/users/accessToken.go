package users

import (
	"context"
	"cos-backend-com/src/common/apierror"
	"cos-backend-com/src/common/auth"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"time"

	"github.com/wujiu2020/strip/sessions"

	"github.com/wujiu2020/strip/utils"
)

const (
	tokenExpiresIn     int64 = 7200
	refreshExpiresIn   int64 = 86400 * 7
	PasswordSaltLength       = 8
)

const (
	TokenTypeRefreshToken = "refreshtoken"
	TokenTypeAccessToken  = "accesstoken"
)

var AccessTokens = &accessTokens{
	models.DefaultConnector,
}

type accessTokens struct{ dbconn.Connector }

func (p *accessTokens) NewToken(ctx context.Context, uid flake.ID, ak, sk string) (token *auth.OAuth2Token, err error) {

	netxtId, err := p.nextId(ctx)
	if err != nil {
		return
	}

	token, err = createOAuthToken(netxtId.String(), ak, sk)
	if err != nil {
		return
	}

	query := `
	INSERT INTO access_tokens
		(id, uid, token, refresh, key, secret)
	VALUES (${id}, ${uid}, ${token}, ${refresh}, ${key}, ${secret})
`
	query, args := util.PgMapQuery(query, map[string]interface{}{
		"{id}":      netxtId,
		"{uid}":     uid,
		"{token}":   token.AccessToken,
		"{refresh}": token.RefreshToken,
		"{key}":     ak,
		"{secret}":  sk,
	})

	err = p.Invoke(ctx, func(db dbconn.Q) error {
		_, er := db.ExecContext(ctx, query, args...)
		return er
	})
	if err != nil {
		return
	}
	return
}

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
