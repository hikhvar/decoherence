package store

import (
	"os"
	"time"
)

// currentDataVersion is the current version of the data scheme
const currentDataVersion dataVersion = 1

type dataVersion int

// FileInfo represents the information recorded of the source. The check command will use this information to
type FileInfo struct {
	Checksum     string      `json:"checksum"`
	RelativePath string      `json:"path"`
	OwnerID      uint32      `json:"owner"`
	GroupID      uint32      `json:"group"`
	FileMode     os.FileMode `json:"file_mode"`
	ModTime      time.Time   `json:"mod_time"`
	Size         int64       `json:"size"`
}

// SetTo sets the values of f to those of other, of they are not the null-values
func (f *FileInfo) SetTo(other FileInfo) {
	if other.RelativePath != "" {
		f.RelativePath = other.RelativePath
	}
	if other.Checksum != "" {
		f.Checksum = other.Checksum
	}
	if other.OwnerID != 0 {
		f.OwnerID = other.OwnerID
	}
	if other.GroupID != 0 {
		f.GroupID = other.GroupID
	}
	if other.FileMode != 0 {
		f.FileMode = other.FileMode
	}
	if !other.ModTime.IsZero() {
		f.ModTime = other.ModTime
	}
	if other.Size != 0 {
		f.Size = other.Size
	}

}

func (f *FileInfo) Equals(other *FileInfo) bool {
	return *f == *other
}

// Meta contains metadata regarding the on stored information
type Meta struct {
	// Version
	Version dataVersion `json:"version"`
	// Checkers is a list of recorded data upfront.
	Checkers []string `json:"checkers"`
}
