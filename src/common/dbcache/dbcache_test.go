package dbcache

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
	"github.com/wujiu2020/strip/utils/helpers"

	"cos-backend-com/src/common/flake"
)

var db *sqlx.DB

type TestData struct {
	ID     flake.ID       `json:"id" db:"id"`
	Name   string         `json:"name" db:"name"`
	Detail types.JSONText `json:"detail" db:"detail"`
}

func TestMain(m *testing.M) {
	if os.Getenv("PG_MASTER") == "" {
		log.Println("dbcache test skipped")
		os.Exit(0)
	}

	var err error
	db, err = sqlx.Connect("postgres", os.Getenv("PG_MASTER"))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ret := m.Run()
	db.Close()
	os.Exit(ret)
}

func Test_dbcache(t *testing.T) {
	db.MustExec(`
CREATE SCHEMA IF NOT EXISTS test;
DROP TABLE IF EXISTS _test_dbcache;
CREATE TABLE _test_dbcache (
    id bigint PRIMARY KEY NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    detail jsonb DEFAULT '{}'::jsonb NOT NULL
);
INSERT INTO _test_dbcache (id, name, detail) VALUES
	(1234567890, 'name','{"x": 62, "y": 51}');
`)

	table := NewTable(TableConfig{
		Table:    "_test_dbcache",
		IdName:   "id",
		Typ:      reflect.TypeOf(TestData{}),
		ExpireIn: 3600,
	})

	dc := &DBCache{log: helpers.X, db: db}
	res, err := dc.fetchDatas(map[*Table][]*Data{
		table: {table.Data(1234567890)},
	})
	assert.NoError(t, err)
	assert.Len(t, res, 1)

	data := res[0].Value.(*TestData)
	assert.NotNil(t, data)
	assert.Equal(t, `{"detail":{"x":62,"y":51},"id":1234567890,"name":"name"}`, string(res[0].bytes))

	d := &struct {
		X, Y int
	}{}

	err = json.Unmarshal(data.Detail, d)
	assert.NoError(t, err)
	assert.Equal(t, 62, d.X)
	assert.Equal(t, 51, d.Y)
}
