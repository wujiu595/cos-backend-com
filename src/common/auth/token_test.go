package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SuInfo(t *testing.T) {

	suInfo, err := CreateSuInfo(1024, "accessToken")
	assert.Equal(t, "aD03OWI0MDU0YzFiNjM4ZmY3MWU0OTM3NjFhY2Y1NDE0NzQzNGY0YzczJnN1aWQ9MTAyNA==", suInfo)
	assert.NoError(t, err)

	suid, err := VerifySuInfo(suInfo, "accessToken")
	assert.NoError(t, err)
	assert.Equal(t, 1024, suid.Int())

	suid, err = VerifySuInfo(suInfo, "accessTokenX")
	assert.Error(t, err)
	assert.Equal(t, 0, suid.Int())
}
