package store

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestFileInfoSetTo(t *testing.T) {
	testcases := []struct {
		name   string
		a      FileInfo
		b      FileInfo
		expect FileInfo
	}{
		{
			name: "merge",
			a: FileInfo{
				Checksum: "foo",
			},
			b: FileInfo{
				Size: 32,
			},
			expect: FileInfo{
				Checksum: "foo",
				Size:     32,
			},
		},
		{
			name: "complete replace",
			a: FileInfo{
				Checksum:     "foo",
				RelativePath: "bar",
				OwnerID:      1,
				GroupID:      2,
				FileMode:     os.ModePerm,
				ModTime:      time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC),
				Size:         3,
			},

			b: FileInfo{
				Checksum:     "bay",
				RelativePath: "chacka",
				OwnerID:      3,
				GroupID:      4,
				FileMode:     0645,
				ModTime:      time.Date(9, 8, 7, 6, 4, 3, 2, time.UTC),
				Size:         5,
			},
			expect: FileInfo{
				Checksum:     "bay",
				RelativePath: "chacka",
				OwnerID:      3,
				GroupID:      4,
				FileMode:     0645,
				ModTime:      time.Date(9, 8, 7, 6, 4, 3, 2, time.UTC),
				Size:         5,
			},
		},
		{
			name: "don't replace with null",
			a: FileInfo{
				Checksum:     "bay",
				RelativePath: "chacka",
				OwnerID:      3,
				GroupID:      4,
				FileMode:     0645,
				ModTime:      time.Date(9, 8, 7, 6, 4, 3, 2, time.UTC),
				Size:         5,
			},
			b: FileInfo{},
			expect: FileInfo{
				Checksum:     "bay",
				RelativePath: "chacka",
				OwnerID:      3,
				GroupID:      4,
				FileMode:     0645,
				ModTime:      time.Date(9, 8, 7, 6, 4, 3, 2, time.UTC),
				Size:         5,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.a.SetTo(tc.b)
			assert.True(t, tc.expect.Equals(&tc.a))
		})
	}
}

func TestFileInfoFields(t *testing.T) {
	fields := []string{"Checksum", "RelativePath", "OwnerID", "GroupID", "FileMode", "ModTime", "Size"}
	f := FileInfo{}
	assert.Equal(t, fields, f.Fields())
}

func TestMustValueOf(t *testing.T) {
	f := FileInfo{Checksum: "foo"}
	assert.Equal(t, "foo", f.MustValueOf("Checksum"))
	defer func() {
		recover()
	}()
	f.MustValueOf("foBar")
	assert.Fail(t, "should never reach this codeline. MustValueOf panics at unknown field.")
}
