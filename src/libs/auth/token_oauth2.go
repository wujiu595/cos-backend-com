package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/wujiu2020/strip"
	"qiniupkg.com/x/rpc.v7"
)

type OAuth2Config struct {
	TokenURL     string
	ClientId     string
	ClientSecret string

	Token     *BearerToken
	Mutex     *sync.RWMutex
	Transport http.RoundTripper
	Log       strip.Logger
}

type oauth2TokenServer struct {
	OAuth2Config
	Client *rpc.Client
}

var _ TokenServer = new(oauth2TokenServer)

func NewOAuth2TokenServer(cfg OAuth2Config) (TokenServer, error) {
	if cfg.Mutex == nil && cfg.Token != nil {
		return nil, fmt.Errorf("mutex cannot be nil when has token store")
	}
	if cfg.Log == nil {
		return nil, fmt.Errorf("log cannot be empty")
	}
	if cfg.Mutex == nil {
		cfg.Mutex = &sync.RWMutex{}
	}
	if cfg.Token == nil {
		cfg.Token = &BearerToken{}
	}
	if cfg.Transport == nil {
		cfg.Transport = http.DefaultTransport
	}
	return &oauth2TokenServer{
		OAuth2Config: cfg,
		Client:       &rpc.Client{&http.Client{Transport: cfg.Transport, Timeout: 10 * time.Second}},
	}, nil
}

func (p *oauth2TokenServer) GetToken() *BearerToken {
	p.Mutex.RLock()
	token := *p.Token
	p.Mutex.RUnlock()
	return &token
}

func (p *oauth2TokenServer) RefreshToken() (newToken *BearerToken, err error) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	token := p.Token

	defer func() {
		if newToken != nil {
			*token = *newToken
		}
	}()

	if token.RefreshToken != "" {
		newToken, err = p.exchangeByRefreshToken(token.RefreshToken)
		if err != nil {
			if v, ok := err.(*rpc.ErrorInfo); ok && v.Code == 401 && (v.Errno == 401 || v.Errno == 20001 || v.Errno == 20004) {
				goto renew
			}
			p.Log.Warn("oauth2TokenServer.exchangeByRefreshToken RefreshToken", err)
			return
		}
		return
	}

renew:
	newToken, err = p.exchangeByClientCredentials(p.ClientId, p.ClientSecret)
	if err != nil {
		p.Log.Warn("oauth2TokenServer.exchangeByClientCredentials RefreshToken", err)
	}
	return
}

func (p *oauth2TokenServer) exchangeByClientCredentials(clientId, clientSecret string) (token *BearerToken, err error) {
	nowUnix := time.Now().Unix()
	param := url.Values{
		"grant_type": {"client_credentials"},
	}

	req, err := http.NewRequest("POST", p.TokenURL+"?"+param.Encode(), nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(clientId, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.Client.Do(nil, req)
	if err != nil {
		return
	}

	var res OAuth2Token
	err = rpc.CallRet(nil, &res, resp)
	if err == nil {
		b := res.ToBearerToken(nowUnix)
		token = &b
	}
	return
}

func (p *oauth2TokenServer) exchangeByRefreshToken(refreshToken string) (token *BearerToken, err error) {
	return exchangeByRefreshToken(p.Client, p.TokenURL, refreshToken)
}

func exchangeByRefreshToken(client *rpc.Client, tokenURL, refreshToken string) (token *BearerToken, err error) {
	nowUnix := time.Now().Unix()
	param := map[string][]string{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}
	var res OAuth2Token
	err = client.CallWithForm(nil, &res, "POST", tokenURL, param)
	if err == nil {
		b := res.ToBearerToken(nowUnix)
		token = &b
	}
	return
}
