package flake

import (
	"errors"
	"sync"
	"time"
)

type SnowFlake struct {
	cfg Config

	shardId  int64
	maxSeqId int64

	lastTimestamp int64
	sequence      int64
	lock          sync.Mutex
}

// init SnowFlake
func NewSnowFlake(shardId int64, cfg Config) (p *SnowFlake, err error) {
	if shardId < 0 || shardId > -1^(-1<<cfg.ShardBits) {
		err = errors.New("invalid shardId not in range")
		return
	}
	p = &SnowFlake{cfg: cfg, shardId: shardId, maxSeqId: -1 ^ (-1 << cfg.SeqBits)}
	return
}

// generate next unique ID
func (p *SnowFlake) Next() ID {
	p.lock.Lock()
	defer p.lock.Unlock()

	ts := timestamp()
	if ts == p.lastTimestamp {
		p.sequence = (p.sequence + 1) & p.maxSeqId
		if p.sequence == 0 {
			ts = p.waitNextMilli(ts)
		}
	} else {
		p.sequence = 0
	}
	p.lastTimestamp = ts

	return PackBits(p.cfg.ShardBits, p.cfg.SeqBits, p.lastTimestamp-p.cfg.Epoch, p.shardId, p.sequence)
}

// sequance exhausted. wait till next millisocond
func (p *SnowFlake) waitNextMilli(ts int64) int64 {
	for ts <= p.lastTimestamp {
		time.Sleep(100 * time.Microsecond)
		ts = timestamp()
	}
	return ts
}

// convert from nanoseconds to milliseconds
func timestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

// pack bits into a snowflake
func PackBits(shardBits, seqBits uint8, msec, shardId, seqId int64) ID {
	return ID(msec<<(shardBits+seqBits) | shardId<<seqBits | seqId)
}
