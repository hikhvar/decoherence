package checker

import (
	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/stretchr/testify/assert"
	"os"
	"os/user"
	"path"
	"strconv"
	"testing"
	"time"
)

const fixturesDir = "../../fixtures"

type fakeFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (f *fakeFileInfo) Name() string {
	return f.name
}
func (f *fakeFileInfo) Size() int64 {
	return f.size
}
func (f *fakeFileInfo) Mode() os.FileMode {
	return f.mode
}
func (f *fakeFileInfo) ModTime() time.Time {
	return f.modTime
}
func (f *fakeFileInfo) IsDir() bool {
	return f.isDir
}
func (f *fakeFileInfo) Sys() interface{} {
	return nil
}

func TestFileInfoExtract(t *testing.T) {
	uid, gid := mustCurrentIDs()
	testcases := []struct {
		name        string
		path        string
		info        os.FileInfo
		expectInfo  store.FileInfo
		expectError string
	}{
		{
			name:        "file does not exists",
			path:        "foo",
			expectInfo:  store.FileInfo{},
			expectError: "could not stat file: no such file or directory",
		},
		{
			name: "isDirectory",
			path: path.Join(fixturesDir, "example-dir"),
			info: &fakeFileInfo{
				name:    "example-dir",
				size:    4,
				mode:    0671,
				modTime: time.Date(2012, 11, 10, 9, 8, 7, 6, time.UTC),
				isDir:   true,
			},
			expectError: "",
			expectInfo: store.FileInfo{
				OwnerID:  uid,
				GroupID:  gid,
				FileMode: 0671,
				Size:     4,
				ModTime:  time.Date(2012, 11, 10, 9, 8, 7, 6, time.UTC),
			},
		},
		{
			name: "normal",
			path: path.Join(fixturesDir, "example.sh"),
			info: &fakeFileInfo{
				name:    "example.sh",
				size:    10,
				mode:    0671,
				modTime: time.Date(2012, 11, 10, 9, 8, 7, 6, time.UTC),
				isDir:   false,
			},
			expectError: "",
			expectInfo: store.FileInfo{
				OwnerID:  uid,
				GroupID:  gid,
				FileMode: 0671,
				Size:     10,
				ModTime:  time.Date(2012, 11, 10, 9, 8, 7, 6, time.UTC),
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			f := FileInfo{}
			info, err := f.Extract(tc.path, tc.info)
			assert.Equal(t, tc.expectInfo, info)
			if tc.expectError == "" {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, tc.expectError, err.Error())
			}
		})
	}
}

func TestFileInfoShouldRunParallel(t *testing.T) {
	f := FileInfo{}
	assert.False(t, f.ShouldRunParallel())
}

func TestFileInfoName(t *testing.T) {
	f := FileInfo{}
	assert.Equal(t, "fileInfo", f.Name())
}

func mustCurrentIDs() (uid uint32, gid uint32) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	gidInt, err := strconv.Atoi(user.Gid)
	if err != nil {
		panic(err)
	}
	uidInt, err := strconv.Atoi(user.Uid)
	if err != nil {
		panic(err)
	}
	return uint32(uidInt), uint32(gidInt)
}
