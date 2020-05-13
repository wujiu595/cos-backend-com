package providers

import (
	"github.com/mediocregopher/radix.v2/pool"

	"cos-backend-com/src/common/sesslimiter"
)

// SessionLimiter returns a new session limiter
func SessionLimiter(redisPool *pool.Pool) *sesslimiter.Limiter {
	return &sesslimiter.Limiter{
		Pool: redisPool,
	}
}
