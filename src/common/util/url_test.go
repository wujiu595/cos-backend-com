package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeURI(t *testing.T) {
	correctResult := `https://www.domain.com/%E4%B8%AD%E6%96%87/path/?name=%E5%90%8D%E5%AD%97&path=/live/&param=%20#%20@`
	raw := "https://www.domain.com/中文/path/?name=名字&path=/live/&param= # @"
	encoded := EncodeURI(raw)
	assert.Equal(t, correctResult, encoded)
}
