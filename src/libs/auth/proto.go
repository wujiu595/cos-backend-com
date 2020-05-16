package auth

import "time"

const (
	LoginToken   = "login_token"
	LoginRefresh = "login_refresh"
	LoginExpired = "login_expired"
)

type OAuth2Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (p *OAuth2Token) ToBearerToken(nowUnix int64) BearerToken {
	return BearerToken{
		AccessToken:  p.AccessToken,
		RefreshToken: p.RefreshToken,
		ExpiresAt:    nowUnix + p.ExpiresIn - 120,
	}
}

type BearerToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"`
}

func (p *BearerToken) IsExpired() bool {
	return time.Now().Unix() > p.ExpiresAt
}

func (p *BearerToken) IsEmpty() bool {
	return p.AccessToken == "" && p.ExpiresAt == 0
}
