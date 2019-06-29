package prospector

import (
	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/schollz/progressbar/v2"
	"os"
)

type progressBarStore struct {
	bar   *progressbar.ProgressBar
	store Store
}

func NewProgressBarStore(max int, s Store) *progressBarStore {
	return &progressBarStore{
		bar:   progressbar.NewOptions(max, progressbar.OptionSetWriter(os.Stderr)),
		store: s,
	}
}

func (p *progressBarStore) Append(finfo ...store.FileInfo) error {
	defer p.bar.Add(len(finfo))
	return p.store.Append(finfo...)
}
