package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PBKDF2Password(t *testing.T) {
	rawPwd := "123456"
	salt := "secret_salt"
	pwd := EncodePBKDF2Password(rawPwd, salt)
	assert.NotEmpty(t, pwd)

	t.Log(pwd)
	assert.False(t, ValidPBKDF2Password(rawPwd, pwd+" "))
	assert.True(t, ValidPBKDF2Password(rawPwd, pwd))
}
