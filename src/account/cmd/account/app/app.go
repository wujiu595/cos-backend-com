package app

import (
	"cos-backend-com/src/account"
	"cos-backend-com/src/account/routers/users"
	"net/http"
	"os"

	"cos-backend-com/src/common/app"
	"cos-backend-com/src/common/providers/session"
	"cos-backend-com/src/common/util"

	"github.com/jmoiron/sqlx"
	"github.com/mediocregopher/radix.v2/pool"
	s "github.com/wujiu2020/strip"
)

const (
	AppName = "account"
)

func AppInit(tea *s.Strip, confPath string, files ...string) *appConfig {
	app := &appConfig{app.New(tea, AppName), account.Env}
	app.ConfigLoad(app.Env, confPath, files...)
	app.ConfigCheck()
	app.ConfigDB()
	app.ConfigProviders()
	app.ConfigFilters()
	app.ConfigRoutes()
	app.ConfigDone()
	return app
}

type appConfig struct {
	app.AppConfig
	Env *account.Config
}

// 配置完毕，做一些运行时检查与初始化
func (p *appConfig) Start() {
}

func (p *appConfig) ConfigCheck() {
}

func (p *appConfig) ConfigDB() *sqlx.DB {
	return p.ConnectDB()
}

func (p *appConfig) ConfigProviders() {

	redisPool, err := pool.NewCustom("tcp",
		p.Env.Redis.Addr,
		p.Env.Redis.PoolSize,
		util.RedisDialWithSecret(p.Env.Redis.Secret),
	)
	if err != nil {
		p.Logger().Error("redis pool.New:", err)
		os.Exit(1)
	}
	p.Provide(redisPool)

	if err := session.ConfigSessions(p.Strip, redisPool, p.Env.Session); err != nil {
		p.Logger().Error("config sessions:", err)
		os.Exit(1)
	}

	var rt http.RoundTripper
	if err := p.Injector().Find(&rt, ""); err != nil {
		panic(err)
	}

	//oauth2TokenURL := p.Env.Service.Auth + "/oauth2/token"
	//
	//tokenServer, err := auth.NewOAuth2TokenServer(auth.OAuth2Config{
	//	TokenURL:     oauth2TokenURL,
	//	ClientId:     p.Env.AdminAccessKey,
	//	ClientSecret: p.Env.AdminAccessSecret,
	//	Token:        &auth.BearerToken{},
	//	Mutex:        &sync.RWMutex{},
	//	Transport:    rt,
	//	Log:          helpers.X,
	//})
	//p.ProvideAs(accountsdk.AdminAuthTransportProvider(tokenServer), (*auth.AdminRoundTripper)(nil))
	if err != nil {
		panic(err)
	}
}

func (p *appConfig) ConfigFilters() {
}

func (p *appConfig) ConfigRoutes() {
	p.Routers(util.VersionRouter())
	p.Routers(
		s.Router("/login",
			s.Post(users.Guest{}).Action("Login"),
		),
	)
}

func (p *appConfig) ConfigDone() {
}
