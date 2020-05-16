package users

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/auth"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/sdk/account"
	"database/sql"
)

const (
	tokenExpiresIn     int64 = 7200
	refreshExpiresIn   int64 = 86400 * 7
	PasswordSaltLength       = 8
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

func (p *accessTokens) RefreshToken(ctx context.Context, m *account.AccessTokensResult) (token *auth.OAuth2Token, err error) {
	accessToken, err := createAccessToken(m.Id.String(), m.Key, m.Secret)
	if err != nil {
		return
	}

	query := `
	UPDATE access_tokens
		SET token = ${newToken}, updated_at = CURRENT_TIMESTAMP
	WHERE id = ${id}
	RETURNING id
`
	query, args := util.PgMapQuery(query, map[string]interface{}{
		"{id}":       m.Id,
		"{newToken}": accessToken,
	})
	err = p.Invoke(ctx, func(db dbconn.Q) error {
		var id flake.ID
		return db.GetContext(ctx, &id, query, args...)
	})
	if err != nil {
		err = apierror.ErrApiAuthCreateFailed
		return
	}

	token = &auth.OAuth2Token{
		AccessToken:  accessToken,
		RefreshToken: m.Refresh,
		ExpiresIn:    tokenExpiresIn,
	}
	return
}

func (p *accessTokens) VerifyToken(ctx context.Context, token string, typ account.TokenType) (m *account.AccessTokensResult, err error) {
	raw, err := verifyTokenOutter(token, typ)
	if err != nil {
		err = apierror.ErrApiAuthInvalidToken
		return
	}
	value, exp, ok := getValuesFromToken(raw)
	if !ok {
		err = apierror.ErrApiAuthInvalidToken
		return
	}

	var id flake.ID
	err = id.UnmarshalText([]byte(value))
	if err != nil {
		return
	}

	model := &account.AccessTokensResult{}
	err = p.FindById(ctx, id, model)
	if err != nil {
		if err == sql.ErrNoRows {
			err = apierror.ErrApiAuthExpiredToken
			return
		}
		return
	}

	err = verifyTokenInner(raw, exp, model.Key, model.Secret)
	if err != nil {
		return
	}

	m = model
	return
}

func (p *accessTokens) FindById(ctx context.Context, id flake.ID, m interface{}) (err error) {
	query := `
	SELECT * FROM access_tokens WHERE id = $1
`
	err = p.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, m, query, id)
	})
	return
}
