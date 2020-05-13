package util

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/util"
)

const (
	DefaultExpiry = 8 * time.Second
	DefaultTries  = 32
	DefaultDelay  = 512 * time.Millisecond
	DefaultFactor = 0.01
)

var (
	ErrMaxTries = errors.New("reached max tries")
)

type RedisMutex struct {
	Name   string
	Expiry time.Duration

	Tries int
	Delay time.Duration

	Factor float64

	value string
	until time.Time

	redis *pool.Pool
	nodem sync.Mutex
}

func NewRedisMutex(name string, redisPool *pool.Pool) *RedisMutex {
	return &RedisMutex{
		Name:   name,
		Expiry: DefaultExpiry,
		Delay:  DefaultDelay,
		Factor: DefaultFactor,
		Tries:  DefaultTries,
		redis:  redisPool,
	}
}

type RedisSync struct {
	redis *pool.Pool
}

func NewRedisSync(redisPool *pool.Pool) *RedisSync {
	return &RedisSync{redisPool}
}

func (r *RedisSync) NewRedisMutex(name string) *RedisMutex {
	return NewRedisMutex(name, r.redis)
}

func (m *RedisMutex) Lock() error {
	m.nodem.Lock()
	defer m.nodem.Unlock()

	return m.lock()
}

func (m *RedisMutex) LockKeeper() (unlocked <-chan struct{}, err error) {
	m.nodem.Lock()
	defer m.nodem.Unlock()

	err = m.lock()
	if err != nil {
		return
	}

	sigCh := make(chan struct{})
	go func() {
		for {
			time.Sleep(m.Expiry / 2)
			if m.touch() {
				continue
			}
			close(sigCh)
			return
		}
	}()
	unlocked = sigCh
	return
}

func (m *RedisMutex) lock() error {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	value := base64.StdEncoding.EncodeToString(b)

	expiry := m.Expiry
	if expiry == 0 {
		expiry = DefaultExpiry
	}

	retries := m.Tries
	if retries == 0 {
		retries = DefaultTries
	}

	delay := m.Delay
	if delay == 0 {
		delay = DefaultDelay
	}

	for i := 0; i < retries; i++ {
		ok, err := m.tryLock(m.Name, value, expiry)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}

		if i == retries-1 {
			break
		}

		time.Sleep(delay)
	}

	return ErrMaxTries
}

func (m *RedisMutex) tryLock(name, value string, expiry time.Duration) (bool, error) {
	start := time.Now()
	conn, err := m.redis.Get()
	if err != nil {
		return false, err
	}
	defer m.redis.Put(conn)

	var ok bool
	resp := conn.Cmd("SET", name, value, "NX", "PX", int(expiry/time.Millisecond))
	if v, _ := resp.Str(); v == "OK" {
		ok = true
	}

	factor := m.Factor
	if factor == 0 {
		factor = DefaultFactor
	}

	until := time.Now().Add(expiry - time.Now().Sub(start) - time.Duration(int64(float64(expiry)*factor)) + 2*time.Millisecond)
	if ok && time.Now().Before(until) {
		m.value = value
		m.until = until
		return true, nil
	}

	util.LuaEval(conn, delScript, 1, name, m.value)
	return false, nil
}

func (m *RedisMutex) Touch() bool {
	m.nodem.Lock()
	defer m.nodem.Unlock()

	return m.touch()
}

func (m *RedisMutex) touch() bool {
	value := m.value
	if value == "" {
		// redis mutex: touch of unlocked mutex
		return false
	}

	expiry := m.Expiry
	if expiry == 0 {
		expiry = DefaultExpiry
	}
	reset := int(expiry / time.Millisecond)

	conn, err := m.redis.Get()
	if err != nil {
		return false
	}
	defer m.redis.Put(conn)

	resp := util.LuaEval(conn, touchScript, 1, m.Name, value, reset)
	if v, err := resp.Str(); err == nil && v != "ERR" {
		return true
	}
	return false
}

func (m *RedisMutex) Unlock() bool {
	m.nodem.Lock()
	defer m.nodem.Unlock()

	value := m.value
	if value == "" {
		panic("redis mutex: unlock of unlocked mutex")
	}

	m.value = ""
	m.until = time.Unix(0, 0)

	conn, err := m.redis.Get()
	if err != nil {
		return false
	}
	defer m.redis.Put(conn)

	resp := util.LuaEval(conn, delScript, 1, m.Name, value)
	if resp.Err == nil {
		return true
	}
	return false
}

var delScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`

var touchScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("set", KEYS[1], ARGV[1], "xx", "px", ARGV[2])
else
	return "ERR"
end`
