package session

import (
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/sessions"
)

type SessionConfig struct {
	Secret       string `conf:"secret"`
	CookiePreifx string `conf:"cookie_prefix"`
	CookieSecure bool   `conf:"cookie_secure"`
	CookieExpire int    `conf:"cookie_expire"`
	CookieDomain string `conf:"cookie_domain"`
	ExpiredIn    int    `conf:"expired_in"`
}

func CreateCookieName(prefix string) string {
	return prefix + "_SESSION"
}

func ConfigSessions(stp *strip.Strip, redisPool *pool.Pool, config SessionConfig) (err error) {
	sessionConfig := sessions.RedisConfig{
		Config: sessions.Config{
			SecretKey:     config.Secret,
			SessionExpire: config.ExpiredIn,
		},
		KeyPrefix: "sess:",
		Client:    redisPool,
	}

	sessionProvider, err := sessions.NewRedisProvider(sessionConfig)
	if err != nil {
		return
	}

	sessionManager := sessions.NewSessionManager(sessionProvider)
	stp.ProvideAs(sessionProvider, (*sessions.SessionProvider)(nil))
	stp.Provide(sessionManager)
	stp.Provide(SessionStore())
	stp.Provide(cookieConfig(config))
	return
}

func cookieConfig(config SessionConfig) interface{} {
	cfg := &sessions.CookieConfig{
		CookieName:         CreateCookieName(config.CookiePreifx),
		CookieRememberName: config.CookiePreifx + "_REMEMBER",
		CookieSecure:       config.CookieSecure,
		CookieDomain:       config.CookieDomain,
		CookieExpire:       0,
		RememberExpire:     3600 * 24 * 10,
	}
	return func() *sessions.CookieConfig {
		return cfg
	}
}
