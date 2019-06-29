package store

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"os"
	"sync"
)

type fileContent struct {
	Meta  Meta       `json:"meta"`
	Files []FileInfo `json:"files"`
}

type JSON struct {
	file    io.WriteCloser
	content *fileContent
	lock    sync.Mutex
}

func (j *JSON) Files() []FileInfo {
	return j.content.Files
}

func NewWriteJSON(path string, meta Meta) (*JSON, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open json datafile")
	}
	return &JSON{
		file: f,
		content: &fileContent{
			Meta: meta,
		}}, nil
}

func NewReadJSON(path string) (*JSON, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open json file")
	}
	defer f.Close()
	j := JSON{}
	return &j, errors.Wrap(json.NewDecoder(f).Decode(&j.content), "failed to decode JSON")
}

func (j *JSON) Append(f ...FileInfo) error {
	j.lock.Lock()
	j.content.Files = append(j.content.Files, f...)
	j.lock.Unlock()
	return nil
}

func (j *JSON) Close() error {
	j.lock.Lock()
	defer j.lock.Unlock()
	err := json.NewEncoder(j.file).Encode(j.content)
	if err != nil {
		return errors.Wrap(err, "failed to encode file content")
	}
	return errors.Wrap(j.file.Close(), "failed to close data file")
}
