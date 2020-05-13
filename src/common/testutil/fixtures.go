package testutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 开启以后，自动把测试生成的请求与结果，写入文件，可以通过 git diff 对比检查
var update = os.Getenv("UPDATE_FIXTURES") != ""

// 设置为 1 以后，强制跳过更新，用于 ci 测试，避免代码里写死 update=true 导致测试始终通过
var enabledUpdate = os.Getenv("UPDATE_FIXTURES_DISABLED") != "1"

type Fixtures struct {
	// 默认从 env 设置以后，可以指定这个进行覆盖
	Update   bool
	Filters  []*FixedFilter
	path     string
	fixtures []*fixture
}

func NewFixtures(path string) *Fixtures {
	return &Fixtures{path: filepath.Join(testsDir, path), Update: update, Filters: append([]*FixedFilter{}, DefaultFixedFilters...)}
}

// Add 增加要测试的结果
func (f *Fixtures) Add(name string, got interface{}, filters ...*FixedFilter) {
	fix := newFixture(name, filepath.Join(f.path, name), got, filters)
	f.fixtures = append(f.fixtures, fix)
}

// Test 执行测试，并且生成测试结果
func (f *Fixtures) Test(t *testing.T) {
	var paths string
	for _, fixture := range f.fixtures {
		idx := 1

		filters := append([]*FixedFilter{}, f.Filters...)
		filters = append(filters, fixture.filters...)

		body := fixture.got
		for _, filter := range filters {
			regex := filter.regex
			body = regex.ReplaceAllFunc(body, func(data []byte) []byte {

				if bytes.Contains(data, []byte(defaultPlaceHolderPrefix)) {
					return data
				}

				ph := defaultPlaceHolderPrefix + fmt.Sprintf("%03d", idx)
				repl := filter.repl(ph)
				idx += 1

				res := regex.ReplaceAll(data, []byte(repl))

				if i := bytes.Index(res, []byte(ph)); i != -1 {
					value := data[i : i+len(data)-(len(res)-len(ph))]
					t.Logf("%s, %s = %s", fixture.name, ph, value)
				}

				return res
			})
		}
		fixture.test(t, body, f.Update)
		paths += fixture.path + "\n"
	}
	fixture := newFixture("", filepath.Join(f.path, "_tests.txt"), paths, nil)
	fixture.test(t, fixture.got, f.Update)
}

type fixture struct {
	name    string
	path    string
	got     []byte
	filters []*FixedFilter
}

func newFixture(name, path string, got interface{}, filters []*FixedFilter) *fixture {
	return &fixture{
		name:    name,
		path:    path,
		got:     []byte(fmt.Sprintf("%s", got)),
		filters: filters,
	}
}

func (f *fixture) test(t *testing.T, got []byte, update bool) {
	if update && enabledUpdate {
		if err := os.MkdirAll(filepath.Dir(f.path), 0755); err != nil {
			t.Error(err)
			return
		} else if err := ioutil.WriteFile(f.path, got, 0644); err != nil {
			t.Error(err)
			return
		}
	}
	want, err := ioutil.ReadFile(f.path)
	if err != nil {
		want = []byte(fmt.Sprintf("%s", err))
	}

	// TODO show diff
	if !assert.Equal(t, string(bytes.TrimSpace(want)), string(bytes.TrimSpace(got))) {
		t.Errorf("%s: does not match, can enable 'Update' flag flush to file\n", f.path)
	}
}
