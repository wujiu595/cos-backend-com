package sesslimiter

import (
	"math"
	"time"

	"github.com/mediocregopher/radix.v2/pool"

	"cos-backend-com/src/common/flake"
)

const (
	prefix = "sesslimiter:"
)

// DomainEnterprise return enterprise domain
func DomainEnterprise(entID flake.ID) string {
	return "ent:" + entID.String()
}

// MemberUser return member key for user
func MemberUser(userID flake.ID) string {
	return "user:" + userID.String()
}

// Limiter activate session limiter
type Limiter struct {
	*pool.Pool
}

// Activate records member in domain
func (l *Limiter) Activate(domain, member string) error {
	return l.ActivateAt(domain, member, time.Now())
}

// ActivateAt records member in domain at specific time
func (l *Limiter) ActivateAt(domain, member string, t time.Time) error {
	cli, err := l.Pool.Get()
	if err != nil {
		return err
	}

	defer l.Pool.Put(cli)

	resp := cli.Cmd("ZADD", prefix+domain, t.Unix(), member)

	return resp.Err
}

// Count returns count of members in the range (now-duration, now]
func (l *Limiter) Count(domain string, duration time.Duration) (int, error) {
	cli, err := l.Pool.Get()
	if err != nil {
		return 0, err
	}

	defer l.Pool.Put(cli)

	now := time.Now()

	resp := cli.Cmd("ZCOUNT", prefix+domain, now.Add(-duration).Unix()+1, now.Unix())
	return resp.Int()
}

// GC removes members in the range [0, now-lifetime]
func (l *Limiter) GC(domain string, lifeTime time.Duration) error {
	cli, err := l.Pool.Get()
	if err != nil {
		return err
	}

	defer l.Pool.Put(cli)

	resp := cli.Cmd("ZREMRANGEBYSCORE", prefix+domain, math.Inf(-1), time.Now().Add(-lifeTime).Unix())
	return resp.Err
}
