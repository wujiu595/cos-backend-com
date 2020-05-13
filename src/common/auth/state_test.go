package auth

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateToken(t *testing.T) {
	now := int64(1494658507)
	secret := "999999"
	value := url.Values{}
	value.Set("u", "12345")

	raw := CreateStateTokenN(secret, value, now)
	assert.Equal(t, "aD01ZGZmZTE2ODU2MTQ5OWEwM2IxMmQyMGRjZDgyMzA4ZiZ0PTE0OTQ2NTg1MDc", raw)

	err := VerifyStateTokenN(secret, value, raw, now, 0)
	assert.NoError(t, err)

	err = VerifyStateTokenN(secret, value, raw, now+10, 11)
	assert.NoError(t, err)

	err = VerifyStateTokenN(secret, value, raw, now+10, 10)
	assert.NoError(t, err)

	err = VerifyStateTokenN(secret, value, raw, now+10, 9)
	assert.Error(t, err)
}
