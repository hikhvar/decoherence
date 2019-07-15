package prospector

import (
	"bytes"
	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testStore struct {
	mem []store.FileInfo
}

func (t *testStore) Append(info ...store.FileInfo) error {
	t.mem = append(t.mem, info...)
	return nil
}

func TestNewProgress(t *testing.T) {
	ts := &testStore{}
	p := NewProgressBarStore(ts)
	buf := bytes.NewBuffer(nil)
	p.output = buf
	p.StartTreeWalk()
	time.Sleep(100 * time.Millisecond)
	assert.True(t, buf.Len() > 1)
	buf.Reset()
	err := p.Append(store.FileInfo{})
	assert.NotNil(t, err)
	p.EndTreeWalk(10)
	assert.True(t, buf.Len() > 1)
	buf.Reset()
	err = p.Append(store.FileInfo{})
	assert.Nil(t, err)
	assert.True(t, buf.Len() > 1)
	buf.Reset()

	p.Finish()
	assert.True(t, buf.Len() > 1)
	buf.Reset()
	assert.True(t, len(ts.mem) == 1)
}
