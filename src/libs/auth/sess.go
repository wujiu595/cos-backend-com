package auth

import (
	. "context"
	"cos-backend-com/src/common/flake"
	"net/http"
	"strings"
	"time"

	"github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/sessions"

	"pkg.jimu.io/account/models"
	"pkg.jimu.io/account/proto"
)

type SignHelper struct {
	Ctx     strip.Context            `inject`
	Req     *http.Request            `inject`
	Rw      http.ResponseWriter      `inject`
	Manager *sessions.SessionManager `inject`
	Sess    sessions.SessionStore    `inject`
	Config  *sessions.CookieConfig   `inject`
}

func (p *SignHelper) SigninUser(ctx Context, uid flake.ID, publicSecret, privateSecret string) (sess sessions.SessionStore, err error) {
	oauthToken, err := models.AccessTokens.NewToken(ctx, uid, publicSecret, privateSecret)
	if err != nil {
		return
	}
	sess, err = p.SigninSession(oauthToken, time.Unix(time.Now().Unix()+oauthToken.ExpiresIn, 0))
	return
}

func (p *SignHelper) SigninSession(token *proto.OAuth2Token, expiredAt time.Time) (sess sessions.SessionStore, err error) {
	p.removeOtherSetCookies()
	sess, _, err = p.Manager.Regenerate(p.Config, p.Rw, p.Req)
	if err != nil {
		return
	}
	p.Ctx.ProvideAs(sess, (*sessions.SessionStore)(nil))

	sess.Set(proto.LoginToken, token.AccessToken)
	sess.Set(proto.LoginRefresh, token.RefreshToken)
	sess.Set(proto.LoginExpired, expiredAt.Unix())
	err = sess.Flush()
	return
}

func (p *SignHelper) SigninMinappSession(sess sessions.SessionStore, openId, sessionKey string) {
	sess.Set(proto.LoginWechatMinappOpenId, openId)
	sess.Set(proto.LoginWechatMinappSess, sessionKey)
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
