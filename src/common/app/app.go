package app

import (
	"database/sql"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/gchaincl/sqlhooks"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/utils/helpers"

	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/filters"
	"cos-backend-com/src/common/proto"
	"cos-backend-com/src/common/providers"
	"cos-backend-com/src/common/util"
)

func init() {
	sql.Register("postgresWithHooks", sqlhooks.Wrap(&pq.Driver{}, &dbconn.Hooks{}))
}

type AppConfig struct {
	Name string
	*strip.Strip
	env interface{}

	SignalReceiver func(os.Signal)
}

func New(sp *strip.Strip, name string) AppConfig {
	sigch := make(chan os.Signal)
	signal.Notify(sigch,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	app := AppConfig{name, sp, nil, nil}
	app.SignalReceiver = app.defaultSignalReceiver

	go func() {
		for {
			app.SignalReceiver(<-sigch)
		}
	}()
	return app
}

func (p *AppConfig) ConfigLoad(env interface{}, confPath string, files ...string) {
	p.env = env
	p.Provide(&proto.AppEnv{Env: reflect.Indirect(reflect.ValueOf(env))})

	helpers.LoadClassicEnv(p.Strip, p.Name, env, confPath, files...)
	helpers.UseGlobalLogger(p.Strip)

	filterOptions := make([]interface{}, 0, 2)

	helpers.X.Recover(func() {
		val := reflect.Indirect(reflect.ValueOf(env))
		if elm, ok := util.FindStructElemRecursive(val, reflect.TypeOf(proto.JaegerConfig{})); ok {
			cfg := elm.Interface().(proto.JaegerConfig)

			p.Provide(providers.JaegerTracing(p.Name, cfg))
			filterOptions = append(filterOptions, helpers.HookFunc(func(sp *strip.Strip, flag helpers.HookFlag) {
				if flag != helpers.HookFlagBeforeAll {
					return
				}
				sp.Filter(filters.OpenTracingFilter())
			}))
		}
	})

	filterOptions = append(filterOptions, util.DefaultLoggerOption(p.Strip))

	helpers.LoadClassicProviders(p.Strip)
	helpers.LoadClassicFilters(p.Strip, filterOptions...)

	providers.LoadClassic(p.Strip)

	// global filters, providers
}

func (p *AppConfig) ConnectDB() *sqlx.DB {
	var cfg proto.PGConfig

	val := reflect.Indirect(reflect.ValueOf(p.env))
	if elm, ok := util.FindStructElemRecursive(val, reflect.TypeOf(proto.PGConfig{})); !ok {
		p.Logger().Error("not found postgres config in env")
		os.Exit(1)
	} else {
		cfg = elm.Interface().(proto.PGConfig)
	}

	conn, err := sqlx.Connect("postgresWithHooks", cfg.Master)
	if err != nil {
		p.Logger().Error(err)
		os.Exit(1)
	}
	conn = conn.Unsafe()

	conn.SetMaxOpenConns(cfg.MaxOpen)
	conn.SetMaxIdleConns(cfg.MaxIdle)

	dbconn.RegisterDB(dbconn.Connectors.Default, conn)
	p.Provide(conn)

	return conn
}

func (p *AppConfig) defaultSignalReceiver(sig os.Signal) {
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
		os.Exit(0)
	}
}
