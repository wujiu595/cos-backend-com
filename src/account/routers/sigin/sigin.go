package sigin

/*
import (
	. "context"
	"cos-backend-com/src/common/auth"
	"cos-backend-com/src/libs/sdk/account"
	"net/http"
	"strings"
	"time"

	"github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/sessions"
)

type SignHelper struct {
	Ctx     strip.Context            `inject`
	Req     *http.Request            `inject`
	Rw      http.ResponseWriter      `inject`
	Manager *sessions.SessionManager `inject`
	Sess    sessions.SessionStore    `inject`
	Config  *sessions.CookieConfig   `inject`
}

func (p *SignHelper) SigninUser(ctx Context, user *account.LoginUserResult, publicSecret, privateSecret string) (sess sessions.SessionStore, err error) {
	tokenInfo := &account.TokenInfo{
		Uid: user.Id,
	}
	oauthToken, err := accountmodels.AccessTokens.NewToken(ctx, publicSecret, tokenInfo)
	if err != nil {
		return
	}
	sess, err = p.SigninSession(oauthToken, time.Unix(time.Now().Unix()+oauthToken.ExpiresIn, 0))
	if err != nil {
		return
	}
	return
}

func (p *SignHelper) SigninSession(token *sdkauth.OAuth2Token, expiredAt time.Time) (sess sessions.SessionStore, err error) {
	p.removeOtherSetCookies()
	sess, _, err = p.Manager.Regenerate(p.Config, p.Rw, p.Req)
	if err != nil {
		return
	}
	p.Ctx.ProvideAs(sess, (*sessions.SessionStore)(nil))

	sess.Set(auth.LoginToken, token.AccessToken)
	sess.Set(auth.LoginRefresh, token.RefreshToken)
	sess.Set(auth.LoginExpired, expiredAt.Unix())
	err = sess.Flush()
	return
}

func (p *SignHelper) SigninMinappSession(sess sessions.SessionStore, openId, sessionKey string) {
	sess.Set(auth.LoginWechatMinappOpenId, openId)
	sess.Set(auth.LoginWechatMinappSess, sessionKey)
}

func (p *SignHelper) Signout() {
	p.Sess.Destroy()
	p.removeOtherSetCookies()
	p.Manager.Destroy(p.Config, p.Rw, p.Req)
	// TODO
	// remove token from AccessTokens
}

func (p *SignHelper) removeOtherSetCookies() {
	h := p.Rw.Header()
	values := make([]string, 0, len(h))
	for _, v := range h["Set-Cookie"] {
		if strings.Contains(v, p.Config.CookieName+"=") {
			continue
		}
		values = append(values, v)
	}
	h["Set-Cookie"] = values
}
*/
