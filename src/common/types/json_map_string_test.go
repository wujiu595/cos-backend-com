package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_JSONMapString(t *testing.T) {
	var m JSONMapString

	b, err := m.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `null`, string(b))

	err = m.UnmarshalJSON(nil)
	assert.Nil(t, m)

	err = m.UnmarshalJSON([]byte(`[]`))
	assert.Error(t, err)

	err = m.UnmarshalJSON([]byte(``))
	assert.Error(t, err)

	err = m.UnmarshalJSON([]byte(`{}`))
	assert.NoError(t, err)
	assert.NotNil(t, m)

	b, err = m.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `{}`, string(b))

	v, err := m.Value()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{}`), v)

	m = nil
	err = m.Scan(nil)
	assert.NoError(t, err)
	assert.Nil(t, m)

	err = m.Scan("")
	assert.NoError(t, err)
	assert.Nil(t, m)

	err = m.Scan([]byte{})
	assert.NoError(t, err)
	assert.Nil(t, m)

	err = m.Scan(1)
	assert.Error(t, err)

	err = m.Scan(`{}`)
	assert.NoError(t, err)
	assert.NotNil(t, m)
}
