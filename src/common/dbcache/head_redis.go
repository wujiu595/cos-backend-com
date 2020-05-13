package dbcache

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/util"
)

type redisHead struct {
	client *pool.Pool
}

func NewRedisHeader(client *pool.Pool) Header {
	return &redisHead{client}
}

func (p *redisHead) MultiGet(keys []string) (d [][]byte, err error) {
	conn, err := p.client.Get()
	if err != nil {
		return
	}
	defer p.client.Put(conn)

	args := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		args = append(args, key)
	}

	resp := conn.Cmd("MGET", args...)
	d, err = resp.ListBytes()
	return
}

func (p *redisHead) MultiSet(keys []string, values [][]byte, expires []int) (err error) {
	conn, err := p.client.Get()
	if err != nil {
		return
	}
	defer p.client.Put(conn)

	args := make([]interface{}, len(keys)*3)
	for i, key := range keys {
		args[i] = key
		args[len(keys)+i] = values[i]
		args[len(keys)*2+i] = int(getTimeoutDur(expires[i]).Seconds())
	}

	script := fmt.Sprintf(`
		for i, key in ipairs(KEYS) do
			if ARGV[%d+i] == '0' then
				redis.call('SET', key, ARGV[i])
			else
				redis.call('SET', key, ARGV[i], 'EX', ARGV[%d+i])
			end
		end
	`, len(keys), len(keys))
	resp := util.LuaEval(conn, script, len(keys), args...)
	err = resp.Err
	return
}

// Deprecated
func (p *redisHead) StreamGet(key string) (rc io.ReadCloser, err error) {
	conn, err := p.client.Get()
	if err != nil {
		return
	}
	defer p.client.Put(conn)

	resp := conn.Cmd("GET", key)
	b, err := resp.Bytes()
	if err != nil {
		return
	}
	rc = ioutil.NopCloser(bytes.NewReader(b))
	return
}

// Deprecated
func (p *redisHead) StreamSet(key string, rc io.ReadCloser, expires ...int) (err error) {
	conn, err := p.client.Get()
	if err != nil {
		return
	}
	defer p.client.Put(conn)

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return
	}

	expire := getTimeoutDur(expires...)
	args := []interface{}{key, b}
	if expire > 0 {
		args = append(args, "EX", expire.Seconds())
	}

	resp := conn.Cmd("SET", args...)
	err = resp.Err
	return
}
