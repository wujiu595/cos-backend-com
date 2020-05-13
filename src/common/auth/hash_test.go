package auth

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashValue(t *testing.T) {
	now := int64(1494658507)
	secret := "999999"
	value := map[string]string{
		"u": "12345",
	}
	raw := CreateHashValue(secret, value, now)
	assert.Equal(t, raw, "h=5d8f47c4e2&t=1494658507&u=12345")

	v, now, er := ParseHashValue(secret, raw, []string{"u"})
	assert.NoError(t, er)
	assert.Equal(t, value, v)

	v, now, er = ParseHashValue(secret, "t=1494658507&u=12345&h=5d8f47c4e2", []string{"u"})
	assert.NoError(t, er)
	assert.Equal(t, value, v)
}

func TestHashValues(t *testing.T) {
	now := int64(1494658507)
	secret := "999999"
	value := url.Values{}
	value.Set("u", "12345")

	raw, _ := CreateHashValues(secret, value, now)
	assert.Equal(t, raw, "_h=5dffe168561499a03b12d20dcd82308f&_t=1494658507&u=12345")

	v, now, er := ParseHashValues(secret, raw)
	assert.NoError(t, er)
	assert.Equal(t, value, v)

	v, now, er = ParseHashValues(secret, "_h=5dffe168561499a03b12d20dcd82308f&_t=1494658507&u=12345")
	assert.NoError(t, er)
	assert.Equal(t, value.Encode(), v.Encode())
}

func TestExpiredHash(t *testing.T) {
	now := int64(1494658507)
	secret := "999999"
	value := url.Values{}
	value.Set("u", "12345")

	raw, err := CreateExpiredHash(secret, value, now)
	assert.NoError(t, err)
	assert.Equal(t, "h=5dffe168561499a03b12d20dcd82308f&t=1494658507", raw)

	remain, err := VerifyExpiredHash(secret, value, raw, now+1)
	assert.NoError(t, err)
	assert.Equal(t, 1, remain)

	remain, err = VerifyExpiredHash(secret, value, "wrong raw value", now-10)
	assert.Error(t, err)
	assert.NotEqual(t, 10, remain)

	value.Set("_", "v")
	remain, err = VerifyExpiredHash(secret, value, raw, now-10)
	assert.Error(t, err)
	assert.NotEqual(t, 10, remain)
}
