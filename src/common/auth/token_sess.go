package auth

import (
	"net/http"
	"sync"

	"github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/sessions"
	"qiniupkg.com/x/rpc.v7"

	"cos-backend-com/src/common/apierror"
)

type SessionTokenConfig struct {
	TokenURL  string
	Transport http.RoundTripper
	Log       strip.Logger
	Sess      sessions.SessionStore
}

type sessTokenServer struct {
	SessionTokenConfig
	mutex  sync.RWMutex
	token  *BearerToken
	client *rpc.Client
}

var _ TokenServer = new(sessTokenServer)

func NewSessTokenServer(config SessionTokenConfig) TokenServer {

	var token BearerToken
	sess := config.Sess
	if sess != nil {
		accessToken := sess.Get(LoginToken).String()
		refreshToken := sess.Get(LoginRefresh).String()
		expiredAt := sess.Get(LoginExpired).MustInt64()
		token.AccessToken = accessToken
		token.RefreshToken = refreshToken
		token.ExpiresAt = expiredAt
	}
	return &sessTokenServer{SessionTokenConfig: config, token: &token, client: &rpc.Client{&http.Client{Transport: config.Transport}}}
}

func (p *sessTokenServer) GetToken() *BearerToken {
	p.mutex.RLock()
	token := *p.token
	p.mutex.RUnlock()
	return &token
}

func (p *sessTokenServer) RefreshToken() (*BearerToken, error) {
	p.mutex.RLock()
	tok := *p.token
	p.mutex.RUnlock()

	if tok.IsEmpty() {
		return nil, apierror.ErrNeedLogin
	}

	if tok.RefreshToken == "" {
		return nil, apierror.ErrApiAuthInvalidToken
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if tok.AccessToken != p.token.AccessToken {
		token := *p.token
		return &token, nil
	}

	token, err := exchangeByRefreshToken(p.client, p.TokenURL, tok.RefreshToken)
	if err != nil {
		return nil, err
	}

	*p.token = *token

	sess := p.Sess
	sess.Set(LoginToken, token.AccessToken)
	sess.Set(LoginRefresh, token.RefreshToken)
	sess.Set(LoginExpired, token.ExpiresAt)
	return token, nil
}
