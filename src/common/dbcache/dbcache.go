package dbcache

import (
	"fmt"
	"reflect"

	"github.com/wujiu2020/strip"

	"github.com/jmoiron/sqlx"
)

var logDC = strip.NewReqLogger(strip.DefaultLogPrinter, "dbcache")

type TableConfig struct {
	Table    string
	IdName   string
	Cols     []string
	Typ      reflect.Type
	ExpireIn int
	Prefix   string
}

type Data struct {
	Table   *Table
	IDValue string
	Value   interface{}
	Error   error
	bytes   []byte
}

type Header interface {
	MultiGet(keys []string) ([][]byte, error)
	MultiSet(keys []string, values [][]byte, expires []int) error
}

type Table struct {
	cfg    TableConfig
	prefix string
}

func NewTable(cfg TableConfig) *Table {
	table := &Table{cfg: cfg, prefix: encodePrefix(cfg.Prefix + cfg.Table)}
	logDC.Infof("table: %v, prefix: %v", cfg.Table, table.prefix)
	return table
}

type DBCache struct {
	log  strip.Logger
	db   *sqlx.DB
	head Header
}

// TODO use orm ??
func NewCache(db *sqlx.DB, head Header) *DBCache {
	return &DBCache{
		log:  logDC,
		db:   db,
		head: head,
	}
}

func (p *DBCache) PreloadAll(table *Table) (err error) {
	query := table.buildQuery(0)
	ptr := reflect.New(reflect.SliceOf(table.cfg.Typ))
	err = p.db.Select(ptr.Interface(), query)
	if err != nil {
		return
	}

	elm := ptr.Elem()
	datas := make([]*Data, 0, elm.Len())
	for i := 0; i < elm.Len(); i++ {
		id := fmt.Sprint(p.db.Mapper.FieldByName(elm.Index(i), table.cfg.IdName).Interface())
		data := &Data{Table: table, IDValue: id}
		data.Value = elm.Index(i).Addr().Interface()
		datas = append(datas, data)
	}

	err = p.updateDatasCache(datas)
	return
}

func (p *DBCache) FetchDatas(datas ...*Data) (err error) {
	keys := make([]string, 0, len(datas))
	for _, data := range datas {
		keys = append(keys, data.key())
	}

	res, err := p.head.MultiGet(keys)
	if err != nil {
		return
	}

	missed := make(map[*Table][]*Data)
	for i, b := range res {
		table := datas[i].Table
		if b != nil {
			ptr := reflect.New(table.cfg.Typ)
			inf := ptr.Interface()
			_, er := sqlxMapMarshalToStruct(b, inf)
			if er == nil {
				datas[i].Value = inf
				continue
			}
			p.log.Warn("key:", keys[i], er)
		}
		missed[table] = append(missed[table], datas[i])
	}
	if len(missed) == 0 {
		return
	}

	fixed, err := p.fetchDatas(missed)
	if err != nil || len(fixed) == 0 {
		return
	}

	go func() {
		err := p.updateDatasCache(fixed)
		if err != nil {
			p.log.Warn(err)
		}
	}()
	return
}

func (p *DBCache) FetchDatasForce(datas ...*Data) (err error) {
	datam := make(map[*Table][]*Data)
	for _, data := range datas {
		datam[data.Table] = append(datam[data.Table], data)
	}

	fixed, err := p.fetchDatas(datam)
	if err != nil || len(fixed) == 0 {
		return
	}

	go func() {
		err := p.updateDatasCache(fixed)
		if err != nil {
			p.log.Warn(err)
		}
	}()
	return
}

func (p *DBCache) fetchDatas(datam map[*Table][]*Data) (results []*Data, err error) {
	for table, datas := range datam {
		args := make([]interface{}, 0, len(datas))

		query := table.buildQuery(len(datas))
		byIds := make(map[string]*Data, len(datas))
		for _, data := range datas {
			args = append(args, data.IDValue)
			byIds[data.IDValue] = data
		}

		err = func() (err error) {
			var rows *sqlx.Rows
			rows, err = p.db.Queryx(query, args...)
			if err != nil {
				return
			}
			defer rows.Close()

			row := make(map[string]interface{})
			for rows.Next() {
				err = rows.MapScan(row)
				if err != nil {
					return
				}

				id := fmt.Sprint(row[table.cfg.IdName])
				data := byIds[id]

				var b []byte
				val := reflect.New(table.cfg.Typ)
				inf := val.Interface()
				b, err = sqlxMapMarshalToStruct(row, inf)
				if err != nil {
					data.Error = err
					return
				}
				data.Value = inf
				data.bytes = b

				results = append(results, data)
				delete(byIds, id)
			}
			return
		}()
		if err != nil {
			return
		}

		for _, data := range byIds {
			data.Error = ErrDataNotFound
		}
	}
	return
}

func (p *DBCache) updateDatasCache(datas []*Data) (err error) {
	keys := make([]string, 0, len(datas))
	bs := make([][]byte, 0, len(datas))
	expires := make([]int, 0, len(datas))
	for _, data := range datas {
		keys = append(keys, data.key())
		bs = append(bs, data.bytes)
		expires = append(expires, data.Table.cfg.ExpireIn)
		data.bytes = nil
	}
	err = p.head.MultiSet(keys, bs, expires)
	return
}
