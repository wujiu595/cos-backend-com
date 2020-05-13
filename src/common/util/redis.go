package util

import (
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

func RedisDialWithSecret(secret string) pool.DialFunc {
	return func(network, addr string) (*redis.Client, error) {
		client, err := redis.DialTimeout(network, addr, 3*time.Second)
		if err != nil {
			return nil, err
		}
		if secret != "" {
			if err = client.Cmd("AUTH", secret).Err; err != nil {
				client.Close()
				return nil, err
			}
		}
		client.ReadTimeout = 10 * time.Second
		client.WriteTimeout = 10 * time.Second
		return client, nil
	}
}
