package prospector

import (
	"errors"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/schollz/progressbar/v2"
	"io"
	"os"
	"time"
)

type progressBarStore struct {
	bar     *progressbar.ProgressBar
	spinner *spinner.Spinner
	store   Store
	output  io.Writer
}

func NewProgressBarStore(s Store) *progressBarStore {
	return &progressBarStore{
		store:  s,
		output: os.Stderr,
	}
}

func (p *progressBarStore) Append(finfo ...store.FileInfo) error {
	if p.bar == nil {
		return errors.New("missing previous call to EndTreeWalk to set maximum number of files")
	}
	defer p.bar.Add(len(finfo))
	return p.store.Append(finfo...)
}

func (p *progressBarStore) Finish() {
	err := p.bar.Finish()
	if err != nil {
		fmt.Fprintf(p.output, "failed to finish progress bar: %v\n", err)
	}
}

func (p *progressBarStore) StartTreeWalk() {
	p.spinner = startSpinner(p.output)
}

func (p *progressBarStore) EndTreeWalk(foundFiles int) {
	p.spinner.Stop()
	fmt.Fprintf(p.output, "Found %d files.\n", foundFiles)
	p.bar = progressbar.NewOptions(foundFiles, progressbar.OptionSetWriter(p.output))
}

func startSpinner(out io.Writer) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond) // Build our new spinner
	s.Writer = out
	s.Suffix = " Reading files"
	s.FinalMSG = "Read all files, start parallel\n"
	mustNotError(s.Color("red"))
	s.Start()
	return s
}

func mustNotError(err error) {
	if err != nil {
		panic("This error must not happen since the error condition is checked at compile time.: " + err.Error())
	}
}
