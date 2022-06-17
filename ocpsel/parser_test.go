package ocpsel

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestParseText(t *testing.T) {
	tests := []struct {
		filePath   string
		exNumChunk int
	}{
		{"empty.in", 0},
		{"selelist.in", 2},
	}

	for _, test := range tests {
		tt := test

		in, err := ioutil.ReadFile(filepath.Join("testdata", tt.filePath))
		if err != nil {
			t.Errorf("%s %s", tt.filePath, err)
		}

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got := ParseText(in)
			if len(got) != tt.exNumChunk {
				t.Errorf("test: %s, got: %d chunks (%+v), expect: %d chunks", tt.filePath, len(got), got, tt.exNumChunk)
			}
		})
	}

}
