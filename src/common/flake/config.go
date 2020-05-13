package flake

import "time"

type Config struct {
	// custom epoch - time offset, milliseconds
	Epoch int64

	// number of bit allocated for server id
	ShardBits uint8
	// number of bit allocated for sequence id
	SeqBits uint8
}

var (
	DBEpoch = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	// 保留一份配置在这里用于解析 ID
	DBConfig = Config{
		Epoch:     DBEpoch.UnixNano() / 1e6,
		ShardBits: 10,
		SeqBits:   12,
	}
	// 保留在这里，用于兼容 parser，应该被重构
	DBFlake, _ = NewSnowFlake(1, DBConfig)
)
