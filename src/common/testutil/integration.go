package testutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wujiu2020/strip/inject"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	dsnPgMaster            = os.Getenv("PG_MASTER")
	databaseResetSql       = filepath.Join(os.Getenv("PROJECT_DIR"), "hack/files/database_test_reset.sql")
	databaseInitialSql     = filepath.Join(os.Getenv("PROJECT_DIR"), "hack/files/database.sql")
	databaseDataSql        = filepath.Join(os.Getenv("PROJECT_DIR"), "hack/files/database_test_data.sql")
	databaseChinaRegionSql = filepath.Join(os.Getenv("PROJECT_DIR"), "hack/files/china_regions.sql")
	databaseSwitchSql      = filepath.Join(os.Getenv("PROJECT_DIR"), "hack/files/database_switch.sql")
	testsDir               = "_tests"

	DefaultDropReqHeaders = []string{
		"User-Agent",
		"X-Reqid",
		"X-Powered-By",
		"X-Response-Time",
	}
	DefaultDropRespHeaders = []string{
		"Set-Cookie",
		"Date",
		"X-Reqid",
		"X-Powered-By",
		"X-Response-Time",
		"Content-Length",
	}
)

type LogFatal interface {
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

type IntegrationTestConfig struct {
	DB              *sqlx.DB
	Log             LogFatal
	Handler         http.HandlerFunc
	HttpClient      *http.Client
	DropReqHeaders  []string
	DropRespHeaders []string
	AutoStart       bool
	StartFunc       func()
	Injector        inject.Injector
}

type IntegrationTest struct {
	IntegrationTestConfig

	Server    *httptest.Server
	Clock     *Clock
	startOnce sync.Once
}

func NewIntegrationTest(config IntegrationTestConfig) *IntegrationTest {
	if config.HttpClient == nil {
		config.HttpClient = http.DefaultClient
	}
	config.Log = log.New(os.Stderr, "", log.LstdFlags)
	if config.DropReqHeaders == nil {
		config.DropReqHeaders = DefaultDropReqHeaders
	}
	if config.DropRespHeaders == nil {
		config.DropRespHeaders = DefaultDropRespHeaders
	}
	p := &IntegrationTest{
		IntegrationTestConfig: config,
		Server:                httptest.NewServer(config.Handler),
		Clock:                 NewClock(time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local)),
	}
	if config.AutoStart {
		p.Start()
	}
	return p
}

func (p *IntegrationTest) Start() {
	p.startOnce.Do(func() {
		if p.IntegrationTestConfig.StartFunc != nil {
			p.IntegrationTestConfig.StartFunc()
		}
	})
}

// Close 测试结束时关闭 server 和 db
func (p *IntegrationTest) Close() {
	p.Server.Close()
	p.DB.Close()
}

// ResetDatabase 重置数据库，指定 loadData 可以加载初始数据
func (p *IntegrationTest) ResetDatabase(loadData bool) {
	p.LoadSql(databaseResetSql)
	p.LoadSql(databaseInitialSql)
	p.LoadSql(databaseSwitchSql)
	if loadData {
		p.LoadSql(databaseDataSql)
		//跑test时没必要跑这个
		// p.LoadSql(databaseChinaRegionSql)
	}
}

// LoadSql 载入执行 sql 文件，如果出错则直接退出
func (p *IntegrationTest) LoadSql(path string) {
	if !strings.HasPrefix(path, "/") {
		path = filepath.Join(testsDir, path)
	}
	if data, err := ioutil.ReadFile(path); err != nil {
		p.Log.Fatal(err)
	} else if _, err := p.DB.Exec(string(data)); err != nil {
		p.Log.Fatalf("%s: %s", path, sqlFriendlyError(string(data), err))
	}
}

// LoadSql 载入执行 sql 文件，如果出错则直接退出
func (p *IntegrationTest) LoadTestJson(path string, dest interface{}) {
	if !strings.HasPrefix(path, "/") {
		path = filepath.Join(testsDir, path)
	}
	if data, err := ioutil.ReadFile(path); err != nil {
		p.Log.Fatal(err)
	} else if err := json.Unmarshal(data, dest); err != nil {
		p.Log.Fatalf("%s: %s", path, err)
	}
}

func sqlFriendlyError(sql string, err error) error {
	pgErr, ok := err.(*pq.Error)
	if !ok {
		return err
	}
	var position int
	if _, scanErr := fmt.Sscanf(pgErr.Position, "%d", &position); scanErr != nil {
		return err
	} else if position <= 0 || position >= len(sql) {
		return err
	}
	start := position - 80
	if start < 0 {
		start = 0
	}
	snippet := sql[start:position]
	return fmt.Errorf("%s near: %s", err, snippet)
}
