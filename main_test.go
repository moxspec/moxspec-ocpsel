package main

import (
	"fmt"
	"testing"
)

func TestSplitEventData(t *testing.T) {
	tests := []struct {
		in  uint64
		ex1 byte
		ex2 byte
		ex3 byte
	}{
		{0, 0, 0, 0},
		{0x020408, 2, 4, 8},
		{0x248, 0, 2, 72},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got1, got2, got3 := splitEventData(tt.in)
			if got1 != tt.ex1 || got2 != tt.ex2 || got3 != tt.ex3 {
				t.Errorf("test: %+v, got: (%x, %x, %x), expect: (%x, %x, %x)", tt, got1, got2, got3, tt.ex1, tt.ex2, tt.ex3)
			}
		})
	}

}
