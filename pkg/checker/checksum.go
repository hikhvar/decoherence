package checker

import (
	"crypto/sha256"
	"fmt"
	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/pkg/errors"
	"io"
	"os"
)

type Checksum struct {
}

func (c *Checksum) Extract(path string, info os.FileInfo) (store.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return store.FileInfo{}, errors.Wrap(err, "failed to open file to compute checksum")
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return store.FileInfo{}, errors.Wrap(err, "failed to copy file content to hasher")
	}

	return store.FileInfo{
		Checksum: fmt.Sprintf("%x", h.Sum(nil)),
	}, nil
}

func (c *Checksum) ShouldRunParallel() bool {
	return true
}

func (c *Checksum) Name() string {
	return "checksum256"
}
