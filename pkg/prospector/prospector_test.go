package prospector

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type staticFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (s staticFileInfo) Name() string {
	return s.name
}

func (s staticFileInfo) Size() int64 {
	return s.size
}

func (s staticFileInfo) Mode() os.FileMode {
	return s.mode
}

func (s staticFileInfo) ModTime() time.Time {
	return s.modTime
}

func (s staticFileInfo) IsDir() bool {
	return s.isDir
}

func (s staticFileInfo) Sys() interface{} {
	return nil
}

func TestProspector_walkFunc(t *testing.T) {
	type fields struct {
		store              Store
		rootPath           string
		files              []store.FileInfo
		progressTracker    ProgressTracker
		inplaceExtractors  []Extractor
		parallelExtractors []Extractor
		maxParallel        int
	}
	type args struct {
		path string
		info os.FileInfo
		err  error
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantFiles []store.FileInfo
	}{
		{
			name: "error not nil",
			args: args{
				path: "foo",
				info: nil,
				err:  errors.New("fff"),
			},
			wantErr: false,
		},
		{
			name: "regular",
			fields: fields{
				rootPath:          "foo",
				inplaceExtractors: []Extractor{&idEx},
			},
			args: args{
				path: "foo/bar",
				info: staticFileInfo{},
			},
			wantFiles: []store.FileInfo{{RelativePath: "bar", OwnerID: uint32(idEx), GroupID: uint32(idEx)}},
		},
		{
			name: "failing extractor",
			fields: fields{
				inplaceExtractors: []Extractor{&errEx},
			},
			args: args{
				path: "foo/bar",
				info: staticFileInfo{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Prospector{
				store:              tt.fields.store,
				rootPath:           tt.fields.rootPath,
				files:              tt.fields.files,
				progressTracker:    tt.fields.progressTracker,
				inplaceExtractors:  tt.fields.inplaceExtractors,
				parallelExtractors: tt.fields.parallelExtractors,
				maxParallel:        tt.fields.maxParallel,
			}
			if err := p.walkFunc(tt.args.path, tt.args.info, tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("Prospector.walkFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.EqualValues(t, tt.wantFiles, p.files)

		})
	}
}

func TestNewProspector(t *testing.T) {
	ts := &testStore{}
	type args struct {
		path        string
		maxParallel int
		store       Store
		extractors  []Extractor
	}
	tests := []struct {
		name string
		args args
		want *Prospector
	}{
		{
			name: "no parallel",
			args: args{
				path:        "foo",
				maxParallel: 0,
				store:       ts,
				extractors: []Extractor{
					parallelExtractor(false), parallelExtractor(true),
				},
			},
			want: &Prospector{
				store:           ts,
				rootPath:        "foo",
				files:           nil,
				progressTracker: NewProgressBarStore(),
				inplaceExtractors: []Extractor{
					parallelExtractor(false), parallelExtractor(true),
				},
				maxParallel: 0,
			},
		},
		{
			name: "with parallel",
			args: args{
				path:        "foo",
				maxParallel: 12,
				store:       ts,
				extractors: []Extractor{
					parallelExtractor(false), parallelExtractor(true),
				},
			},
			want: &Prospector{
				store:           ts,
				rootPath:        "foo",
				files:           nil,
				progressTracker: NewProgressBarStore(),
				inplaceExtractors: []Extractor{
					parallelExtractor(false),
				},
				parallelExtractors: []Extractor{
					parallelExtractor(true),
				},
				maxParallel: 12,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProspector(tt.args.path, tt.args.maxParallel, tt.args.store, tt.args.extractors); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProspector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProspector_Run(t *testing.T) {
	type fields struct {
		store              Store
		rootPath           string
		files              []store.FileInfo
		progressTracker    ProgressTracker
		inplaceExtractors  []Extractor
		parallelExtractors []Extractor
		maxParallel        int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Prospector{
				store:              tt.fields.store,
				rootPath:           tt.fields.rootPath,
				files:              tt.fields.files,
				progressTracker:    tt.fields.progressTracker,
				inplaceExtractors:  tt.fields.inplaceExtractors,
				parallelExtractors: tt.fields.parallelExtractors,
				maxParallel:        tt.fields.maxParallel,
			}
			if err := p.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Prospector.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_worker(t *testing.T) {
	type args struct {
		done       *sync.WaitGroup
		s          Store
		backlog    chan store.FileInfo
		errorChan  chan error
		extractors []Extractor
		rootPath   string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker(tt.args.done, tt.args.s, tt.args.backlog, tt.args.errorChan, tt.args.extractors, tt.args.rootPath)
		})
	}
}
