package prospector

import (
	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

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
