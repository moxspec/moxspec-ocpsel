package ocpsel

import (
	"fmt"
	"math"
	"strings"
)

func splitEventData(ed uint64) (byte, byte, byte) {
	ed1 := byte((ed >> 16) & 0xFF)
	ed2 := byte((ed >> 8) & 0xFF)
	ed3 := byte(ed & 0xFF)
	return ed1, ed2, ed3
}

func pick(ed byte, from, to uint) (byte, error) {
	if from > to || to > 7 {
		return 0, fmt.Errorf("invalid range given")
	}

	// just single-bit
	if from == to {
		return (ed >> from) & 0x01, nil
	}

	var mask byte
	for i := 0; i <= int(to-from); i++ {
		mask = mask + byte(math.Pow(2, float64(i)))
	}

	return (ed >> from) & mask, nil
}

func mustPick(ed byte, from, to uint) byte {
	ret, _ := pick(ed, from, to)
	return ret
}

func shouldBe(name string, ed, expect byte) string {
	if ed != expect {
		return fmt.Sprintf("%s should be 0x%02X but 0x%02X", name, expect, ed)
	}
	return ""
}

func ed1ShouldBe(ed, expect byte) string {
	return shouldBe("ed1", ed, expect)
}

func ed2ShouldBe(ed, expect byte) string {
	return shouldBe("ed1", ed, expect)
}

func ed3ShouldBe(ed, expect byte) string {
	return shouldBe("ed1", ed, expect)
}

func dumpBinary(in interface{}) string {
	var mlen byte
	var tgt uint64
	switch in.(type) {
	case uint8:
		mlen = 8
		tgt = uint64(in.(byte))
	case uint16:
		mlen = 16
		tgt = uint64(in.(uint16))
	case uint32:
		mlen = 32
		tgt = uint64(in.(uint32))
	case uint64:
		mlen = 64
		tgt = in.(uint64)
	case uint:
		mlen = 64
		tgt = uint64(in.(uint))
	default:
		return ""
	}

	amlen := mlen + (mlen/4 - 1) // actual maximum length
	res := make([]string, amlen, amlen)
	var i byte
	var j = amlen - 1
	for i < mlen {
		r := "1"
		if tgt&(1<<i) == 0 {
			r = "0"
		}
		res[j] = r
		if j != 0 && j != amlen-1 && (i+1)%4 == 0 {
			j--
			res[j] = " "
		}
		j--
		i++
	}

	return strings.Join(res, "")
}
