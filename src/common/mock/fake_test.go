package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateFakeMacaddr(t *testing.T) {
	name := "dev-ipc-01"
	macaddr := FakeMacaddr(name)
	assert.Equal(t, "b3:b9:52:13:cb:3a", macaddr)
}

func Test_CreateFakeId(t *testing.T) {
	plantId := FakeId("plant")
	assert.Equal(t, int64(681299897146), plantId.Int64())

	plantId = FakeId("plant-xxxx")
	assert.Equal(t, int64(937682950854), plantId.Int64())

	plantId1 := FakeId("plant-1")
	assert.Equal(t, int64(44649670059360257), plantId1.Int64())

	plantId21 := FakeId("plant-21")
	assert.Equal(t, int64(44649670059360277), plantId21.Int64())

	plantIdMaxSeq := FakeId("plant-65535")
	assert.Equal(t, int64(44649670059425791), plantIdMaxSeq.Int64())

	assert.Panics(t, func() {
		_ = FakeId("plant-65536")
	}, "must overflow max seq id")

	assert.Condition(t, func() bool {
		return plantId1 < plantId21
	})
}
