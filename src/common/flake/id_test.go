package flake

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Snowflake_ID(t *testing.T) {
	var (
		shardBits uint8 = 12
		seqBits   uint8 = 11
		shardId   int64 = (1 << shardBits) - 1
		seqId     int64 = (1 << seqBits) - 1
		msec      int64 = (1<<(63-(shardBits+seqBits)) - 1)
		maxId     int64 = 9223372036854775807
		// maxId     int64  = 14802775804217345
		maxIdBin []byte = make([]byte, 8)
	)
	binary.BigEndian.PutUint64(maxIdBin, uint64(maxId))

	id := PackBits(shardBits, seqBits, msec, shardId, seqId)

	assert.Equal(t, maxId, id.Int64())
	assert.Equal(t, int(maxId), id.Int())
	assert.Equal(t, maxIdBin, id.Bytes())
	assert.Equal(t, fmt.Sprintf(`%d`, maxId), id.String())

	b, err := id.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(`"%d"`, maxId), string(b))

	b, err = id.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(`%d`, maxId), string(b))

	b, err = id.MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, maxIdBin, b)

	value, err := id.Value()
	assert.NoError(t, err)
	assert.Equal(t, maxId, value)

	id = ID(0)
	err = id.UnmarshalJSON([]byte(fmt.Sprintf(`%d`, maxId)))
	assert.NoError(t, err)
	assert.Equal(t, maxId, id.Int64())

	id = ID(0)
	err = id.UnmarshalJSON([]byte(fmt.Sprintf(`"%d"`, maxId)))
	assert.NoError(t, err)
	assert.Equal(t, maxId, id.Int64())

	id = ID(0)
	err = id.UnmarshalText([]byte(fmt.Sprintf(`%d`, maxId)))
	assert.NoError(t, err)
	assert.Equal(t, maxId, id.Int64())

	id = ID(0)
	err = id.UnmarshalBinary(maxIdBin)
	assert.NoError(t, err)
	assert.Equal(t, maxId, id.Int64())

	id = ID(0)
	err = id.Scan(maxIdBin)
	assert.NoError(t, err)
	assert.Equal(t, maxId, id.Int64())

	id = ID(0)
	err = id.Scan(maxIdBin)
	assert.NoError(t, err)
	assert.Equal(t, maxId, id.Int64())

	id = ID(0)
	err = id.Scan([]byte(strconv.FormatInt(maxId, 10)))
	assert.NoError(t, err)
	assert.Equal(t, maxId, id.Int64())

	id = ID(0)
	err = id.Scan(strconv.FormatInt(maxId, 10))
	assert.NoError(t, err)
	assert.Equal(t, maxId, id.Int64())

	id, err = FromString(strconv.FormatInt(maxId, 10))
	assert.NoError(t, err)
	assert.Equal(t, maxId, id.Int64())
}
