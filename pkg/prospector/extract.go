package prospector

import (
	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

// ExtractFileInfo applies the given extractors on the file specified by path.
// root is the starting directory for the record run.
// Already known facts about the file are given by info.
// If either any extractor returns an error, or path is not relative to root, an error is returned
// In case of an error the returned store.FileInfo may be inconsistent
func ExtractFileInfo(extractors []Extractor, root, path string, info os.FileInfo) (store.FileInfo, error) {
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return store.FileInfo{}, errors.Wrapf(err, "relative path from root to file %s does not exists", path)
	}
	f := store.FileInfo{
		RelativePath: relPath,
	}
	for _, ex := range extractors {
		exInfo, err := ex.Extract(path, info)
		if err != nil {
			return store.FileInfo{}, errors.Wrapf(err, "failed to extract information via extractor %s from %s", ex.Name(), path)
		}
		f.SetTo(exInfo)
	}
	return f, nil
}
