package dbcache

import (
	"time"
)

const (
	defaultMemHeaderSize   = 2000
	defaultMemMaxValueSize = 256 * 1024 // 256k, 超过大小的跳过
)

type memHead struct {
	cache *lru.ARCCache
}

type memCacheItem struct {
	expired int64
	value   []byte
}

func NewMemHeader(size int) Header {
	if size <= 0 {
		size = defaultMemHeaderSize
	}
	cache, _ := lru.NewARC(size)
	return &memHead{cache: cache}
}

func (p *memHead) MultiGet(keys []string) (values [][]byte, err error) {
	now := time.Now().Unix()
	values = make([][]byte, len(keys))
	for i, key := range keys {
		if v, ok := p.cache.Get(key); ok {
			item := v.(memCacheItem)
			if item.expired <= 0 || item.expired > now {
				values[i] = item.value
			}
		}
	}
	return
}

func (p *memHead) MultiSet(keys []string, values [][]byte, expires []int) (err error) {
	now := time.Now()
	for i, key := range keys {
		timeout := getTimeoutDur(expires...)
		var expired int64
		if timeout > 0 {
			expired = now.Add(timeout).Unix()
		}
		item := memCacheItem{expired, values[i]}
		p.cache.Add(key, item)
	}
	return
}
