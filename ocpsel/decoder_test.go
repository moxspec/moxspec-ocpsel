package ocpsel

import (
	"fmt"
	"testing"
)

func TestGetDecoder(t *testing.T) {
	tests := []struct {
		sen        uint64
		gen        uint64
		wantsError bool
	}{
		{0, 0x248, true},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			_, _, err := GetDecoder(tt.sen, tt.gen)

			if tt.wantsError && err == nil {
				t.Errorf("test: %+v, wants error but was nil", tt)
			}

			if !tt.wantsError && err != nil {
				t.Errorf("test: %+v, doesn't want error but was %s", tt, err)
			}
		})
	}

}
