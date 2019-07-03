package checker

import (
	"github.com/hikhvar/decoherence/pkg/store"
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func TestChecksumExtract(t *testing.T) {
	testcases := []struct {
		name          string
		path          string
		expectedInfo  store.FileInfo
		expectedError string
	}{
		{
			name: "normal",
			path: path.Join(fixturesDir, "example.sh"),
			expectedInfo: store.FileInfo{
				Checksum: "d26ef5a4b56d86a9319ff6604a31d0dd0d2abadd77124a76dfe7df8f08132fae",
			},
			expectedError: "",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := Checksum{}
			info, err := c.Extract(tc.path, nil)
			assert.Equal(t, tc.expectedInfo, info)
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestChecksumShouldRunParallel(t *testing.T) {
	c := Checksum{}
	assert.True(t, c.ShouldRunParallel())
}

func TestChecksumName(t *testing.T) {
	c := Checksum{}
	assert.Equal(t, "checksum256", c.Name())
}
