package oauth

import (
	"cos-backend-com/src/account/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/users"
	"cos-backend-com/src/libs/sdk/account"

	"github.com/wujiu2020/strip/utils/apires"
)

type Token struct {
	routers.Base
}

type tokenInput struct {
	GrantType    string `param:"grant_type"`
	RefreshToken string `param:"refresh_token"`
}

func (p *Token) GrantToken() (res interface{}) {

	var input tokenInput
	p.Params.BindValuesToStruct(&input)

	switch input.GrantType {
	case "refresh_token":
		if input.RefreshToken == "" {
			p.Log.Warn("refresh_token empty when refresh")
			res = apierror.ErrApiAuthInvalidToken
			return
		}
		res = p.grantByRefreshToken(input.RefreshToken)
	case "client_credentials":
		ak, sk, ok := p.Req.BasicAuth()
		if !ok {
			p.Log.Warn("client_credentials Req.BasicAuth: check failed")
			res = apierror.ErrApiAuthInvalidHeader
			return
		}
		res = p.grantByClientCredentials(ak, sk)

	default:
		p.Log.Warnf("unsupported grant_type: %q", input.GrantType)
		res = apierror.ErrApiAuthInvalidType
	}

	return
}

func (p *Token) grantByRefreshToken(refreshToken string) (res interface{}) {
	tokenInfo, err := users.AccessTokens.VerifyToken(p.Ctx, refreshToken, account.TokenTypeRefreshToken)
	if err != nil {
		p.Log.Warn("accountmodels.AccessTokens.VerifyToken:", err)
		res = apierror.HandleError(err)
		return
	}

	p.Log.Infof("refresh_token ak: %s, uid: %d", tokenInfo.Key, tokenInfo.UId)

	token, err := users.AccessTokens.RefreshToken(p.Ctx, tokenInfo)
	if err != nil {
		p.Log.Warn("accountmodels.AccessTokens.RefreshToken:", err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(token)
	return
}

func (p *Token) grantByClientCredentials(ak, sk string) (res interface{}) {

	if ak == "" || sk == "" {
		res = apierror.HandleError(apierror.ErrApiAuthInvalidInfo)
		return
	}

	var access account.AccessTokensResult
	err := users.AccessKeys.FindByAccessKey(p.Ctx, ak, &access)
	if err != nil {
		p.Log.Warnf("access key of `%s` FindByAccessKey: %v", ak, err)
		res = apierror.HandleError(apierror.ErrApiAuthInvalidInfo)
		return
	}

	p.Log.Infof("client_credentials ak: %s, uid: %d", access.Key, access.UId)

	if access.Secret != sk {
		p.Log.Warnf("access key of `%s` secret key not matched", ak)
		res = apierror.HandleError(apierror.ErrApiAuthInvalidInfo)
		return
	}

	token, err := users.AccessTokens.NewToken(p.Ctx, access.UId, access.Key, access.Secret)
	if err != nil {
		p.Log.Warn("accountmodels.CreateAccessToken:", err)
		res = apierror.HandleError(apierror.ErrApiAuthCreateFailed)
		return
	}

	res = apires.With(token)
	return
}
