package ocpsel

import (
	"fmt"
	"testing"
)

func TestPick(t *testing.T) {
	tests := []struct {
		ed         byte
		from       uint
		to         uint
		ex         byte
		wantsError bool
	}{
		{0xFF, 0, 1, 0x03, false},
		{0xFF, 0, 0, 0x01, false},
		{0xF0, 3, 4, 0x02, false},
		{0xF0, 4, 7, 0x0F, false},
		{0xFF, 5, 0, 0x01, true},
		{0xFF, 0, 8, 0x01, true},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got, err := pick(tt.ed, tt.from, tt.to)

			if tt.wantsError && err == nil {
				t.Errorf("test: %+v, wants error but was nil", tt)
			}

			if !tt.wantsError && got != tt.ex {
				t.Errorf("test: %+v, got: %d (%s), expect: %d", tt, got, err, tt.ex)
			}
		})
	}

}

func TestDumpBinary(t *testing.T) {
	tests := []struct {
		in interface{}
		ex string
	}{
		{uint(248), "0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 1111 1000"},
		{uint64(248), "0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 0000 1111 1000"},
		{uint32(248), "0000 0000 0000 0000 0000 0000 1111 1000"},
		{uint16(248), "0000 0000 1111 1000"},
		{uint8(248), "1111 1000"},
	}

	for _, test := range tests {
		tt := test

		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			got := dumpBinary(tt.in)
			if got != tt.ex {
				t.Errorf("test: %+v, got: %s, expect: %s", tt, got, tt.ex)
			}
		})
	}

}

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
