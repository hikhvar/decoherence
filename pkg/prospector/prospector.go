package prospector

import (
	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type Store interface {
	Append(info ...store.FileInfo) error
}

type ProgressTracker interface {
	Wrap(s Store) Store
	Finish()
	StartTreeWalk()
	EndTreeWalk(foundFiles int)
}

type Extractor interface {
	Extract(path string, info os.FileInfo) (store.FileInfo, error)
	ShouldRunParallel() bool
	Name() string
}

type Prospector struct {
	store              Store
	rootPath           string
	files              []store.FileInfo
	progressTracker    ProgressTracker
	inplaceExtractors  []Extractor
	parallelExtractors []Extractor
	maxParallel        int
}

func NewProspector(path string, maxParallel int, store Store, extractors []Extractor) *Prospector {
	var inplace, parallel []Extractor
	for _, ex := range extractors {
		if ex.ShouldRunParallel() && maxParallel > 1 {
			parallel = append(parallel, ex)
		} else {
			inplace = append(inplace, ex)
		}
	}
	return &Prospector{
		rootPath:           path,
		progressTracker:    NewProgressBarStore(),
		store:              store,
		maxParallel:        maxParallel,
		inplaceExtractors:  inplace,
		parallelExtractors: parallel,
	}
}

func (p *Prospector) Run() error {
	p.progressTracker.StartTreeWalk()
	err := filepath.Walk(p.rootPath, p.walkFunc)
	if err != nil {
		return errors.Wrap(err, "initial prospector run failed")
	}
	p.progressTracker.EndTreeWalk(len(p.files))
	if len(p.parallelExtractors) < 1 {
		return errors.Wrap(p.store.Append(p.files...), "failed to append files")
	}
	p.store = p.progressTracker.Wrap(p.store)
	defer p.progressTracker.Finish()
	var done sync.WaitGroup
	doneChan := make(chan struct{})
	errorChan := make(chan error)
	backlog := make(chan store.FileInfo, p.maxParallel*100)
	for i := 0; i < p.maxParallel; i++ {
		done.Add(1)
		go worker(&done, p.store, backlog, errorChan, p.parallelExtractors, p.rootPath)
	}
	go func() {
		for _, f := range p.files {
			backlog <- f
		}
		close(backlog)
	}()
	go func() {
		done.Wait()
		doneChan <- struct{}{}
	}()
	var errs []error
	for {
		select {
		case err := <-errorChan:
			errs = append(errs, err)
		case <-doneChan:
			return NewMultiError(errs)
		}
	}
}

func (p *Prospector) walkFunc(path string, info os.FileInfo, err error) error {
	if info.IsDir() || !info.Mode().IsRegular() {
		return nil
		// TODO: Figure out how to handle directories and irregular files
	}
	f, err := ExtractFileInfo(p.inplaceExtractors, p.rootPath, path, info)
	if err != nil {
		return errors.Wrap(err, "failed to extract file info")
	}
	p.files = append(p.files, f)
	return nil
}

func worker(done *sync.WaitGroup, s Store, backlog chan store.FileInfo, errorChan chan error, extractors []Extractor, rootPath string) {
	defer done.Done()
	for f := range backlog {
		fPath := path.Join(rootPath, f.RelativePath)
		info, err := os.Stat(fPath)
		if err != nil {
			errorChan <- errors.Wrapf(err, "could not stat %s", fPath)
			continue
		}
		fNew, err := ExtractFileInfo(extractors, rootPath, fPath, info)
		if err != nil {
			errorChan <- errors.Wrapf(err, "could not run parallel extractors on file %s", fPath)
			continue
		}
		f.SetTo(fNew)
		err = s.Append(f)
		if err != nil {
			errorChan <- errors.Wrapf(err, "failed to append file %s to store", fPath)
		}
	}
}
