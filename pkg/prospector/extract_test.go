package prospector

import (
	"errors"
	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type parallelExtractor bool

func (p parallelExtractor) Extract(path string, info os.FileInfo) (store.FileInfo, error) {
	return store.FileInfo{}, nil
}

func (p parallelExtractor) ShouldRunParallel() bool {
	return bool(p)
}

func (p parallelExtractor) Name() string {
	return "parallelExtractor"
}

type errorExtractor string

func (e *errorExtractor) Extract(path string, info os.FileInfo) (store.FileInfo, error) {
	return store.FileInfo{}, errors.New("errormessage")
}

func (e *errorExtractor) ShouldRunParallel() bool {
	return false
}
func (e *errorExtractor) Name() string {
	return "errorExtractor"
}

type idExtractor uint32

func (e *idExtractor) Extract(path string, info os.FileInfo) (store.FileInfo, error) {
	return store.FileInfo{OwnerID: uint32(*e), GroupID: uint32(*e)}, nil
}

func (e *idExtractor) ShouldRunParallel() bool {
	return false
}
func (e *idExtractor) Name() string {
	return "idExtractor"
}

var errEx errorExtractor = "errorExtractor"
var idEx idExtractor = 42

func TestExtractFileInfo(t *testing.T) {

	testcases := []struct {
		name        string
		extractors  []Extractor
		root        string
		path        string
		info        os.FileInfo
		expectInfo  store.FileInfo
		expectError string
	}{
		{
			name:        "not relative path",
			root:        "/foo/bar",
			path:        "./baz",
			expectInfo:  store.FileInfo{},
			expectError: "relative path from root to file ./baz does not exists: Rel: can't make ./baz relative to /foo/bar",
		},
		{
			name:        "no extractor",
			root:        "/foo",
			path:        "/foo/bar",
			expectInfo:  store.FileInfo{RelativePath: "bar"},
			expectError: "",
		},
		{
			name:        "error extractor",
			extractors:  []Extractor{&idEx, &errEx},
			root:        "/foo",
			path:        "/foo/bar",
			expectInfo:  store.FileInfo{},
			expectError: "failed to extract information via extractor errorExtractor from /foo/bar: errormessage",
		},
		{
			name:       "valid extractor",
			extractors: []Extractor{&idEx},
			root:       "/foo",
			path:       "/foo/bar",
			expectInfo: store.FileInfo{
				RelativePath: "bar",
				OwnerID:      uint32(idEx),
				GroupID:      uint32(idEx),
			},
			expectError: "",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := ExtractFileInfo(tc.extractors, tc.root, tc.path, tc.info)
			assert.Equal(t, tc.expectInfo, f)
			if tc.expectError != "" {
				assert.EqualError(t, err, tc.expectError)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
