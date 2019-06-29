package checker

import (
	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"os"
)

type FileInfo struct {
}

func (f *FileInfo) Extract(path string, info os.FileInfo) (store.FileInfo, error) {
	var statT = unix.Stat_t{}
	err := unix.Stat(path, &statT)
	if err != nil {
		return store.FileInfo{}, errors.Wrap(err, "could not stat file")
	}
	return store.FileInfo{
		OwnerID:  statT.Uid,
		GroupID:  statT.Gid,
		FileMode: info.Mode(),
		Size:     info.Size(),
		ModTime:  info.ModTime(),
	}, nil
}

func (f *FileInfo) ShouldRunParallel() bool {
	return false
}

func (f *FileInfo) Name() string {
	return "fileInfo"
}
