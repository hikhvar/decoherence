package prospector

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMultiError(t *testing.T) {
	me := NewMultiError(nil)
	assert.Nil(t, me)
	me = NewMultiError([]error{})
	assert.Nil(t, me)
	me = NewMultiError([]error{errors.New("foo")})
	assert.NotNil(t, me)
}

func TestMultiErrorError(t *testing.T) {
	me := NewMultiError([]error{errors.New("foo"), errors.New("bar")})
	assert.Equal(t, "foo;bar", me.Error())
}
