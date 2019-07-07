package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputeDiffs(t *testing.T) {
	testcases := []struct {
		name string
		old  []FileInfo
		new  []FileInfo
		on   []EqualCriteria
		res  Result
	}{
		{
			name: "some missing, some wrong, some too much",
			old: []FileInfo{
				{
					RelativePath: "shouldBeEqual",
					Checksum:     "equal",
				},
				{
					RelativePath: "onlyInOld",
					Checksum:     "veryOld",
				},
				{
					RelativePath: "wrongInNew",
					Checksum:     "veryOld",
				},
			},
			new: []FileInfo{
				{
					RelativePath: "shouldBeEqual",
					Checksum:     "equal",
				},
				{
					RelativePath: "onlyInNew",
					Checksum:     "veryNew",
				},
				{
					RelativePath: "wrongInNew",
					Checksum:     "veryNew",
				},
			},
			res: Result{
				missingInNew: []diff{
					{
						expected: FileInfo{
							RelativePath: "onlyInOld",
							Checksum:     "veryOld",
						},
					},
				},
				missingInOld: []diff{
					{
						got: FileInfo{
							RelativePath: "onlyInNew",
							Checksum:     "veryNew",
						},
					},
				},
				notEqual: []diff{
					{
						expected: FileInfo{
							RelativePath: "wrongInNew",
							Checksum:     "veryOld",
						},
						got: FileInfo{
							RelativePath: "wrongInNew",
							Checksum:     "veryNew",
						},
					},
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res := ComputeDiffs(tc.old, tc.new, tc.on...)
			assert.Equal(t, tc.res, res)
		})
	}
}
