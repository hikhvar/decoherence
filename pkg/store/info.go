package store

import (
	"os"
	"reflect"
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

// Fields return the names of all fields of the FileInfo struct.
func (f *FileInfo) Fields() (fields []string) {
	s := reflect.ValueOf(f).Elem().Type()
	for i := 0; i < s.NumField(); i++ {
		fields = append(fields, s.Field(i).Name)
	}
	return fields
}

//MustValueOf returns the value of the given fieldName. If the fieldName is not know, this function panics.
func (f *FileInfo) MustValueOf(fieldName string) interface{} {
	return reflect.ValueOf(f).Elem().FieldByName(fieldName).Interface()
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

// EqualCriteria functions compares two FileInfo objects.
type EqualCriteria func(a, b FileInfo) bool

// AllFieldValuesEqual returns true if the values of all struct fields are equal. Be aware, that this will break, if a field
// becomes a pointer type.
func AllFieldValuesEqual(a, b FileInfo) bool {
	return a == b
}

// Equals return true of both FileInfo contains the same information. They do not need the be the same instance.
// The compare criteria can be given via the optional criteria parameter. If there are no criterias given, the AllFieldsEqual criteria is used.
func (f *FileInfo) Equals(other *FileInfo, criteria ...EqualCriteria) bool {
	if len(criteria) == 0 {
		return AllFieldValuesEqual(*f, *other)
	}
	for _, crit := range criteria {
		if !crit(*f, *other) {
			return false
		}
	}
	return true
}

// Meta contains metadata regarding the on stored information
type Meta struct {
	// Version
	Version dataVersion `json:"version"`
	// Checkers is a list of recorded data upfront.
	Checkers []string `json:"checkers"`
}
