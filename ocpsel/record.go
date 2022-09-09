package ocpsel

import (
	"fmt"
	"strings"
)

// record represents a SEL record
type record struct {
	title   string
	genID   uint64
	genName string
	ed1     byte
	r1      string
	ed2     byte
	r2      string
	ed3     byte
	r3      string
}

// Summary returns record summary
func (r record) Summary() string {
	return r.SummaryIndent("")
}

// SummaryIndent returns record summary with given indent
func (r record) SummaryIndent(in string) string {
	var res []string

	res = append(res, fmt.Sprintf("%sDecoded Info:", in))
	if r.title != "" {
		res = append(res, fmt.Sprintf("%s  Summary      : %s", in, r.title))
	}

	if r.r1 != "" {
		res = append(res, fmt.Sprintf("%s  Event Data 1 : 0x%02X, %s ... %s", in, r.ed1, dumpBinary(r.ed1), r.r1))
	} else {
		res = append(res, fmt.Sprintf("%s  Event Data 1 : 0x%02X, %s", in, r.ed1, dumpBinary(r.ed1)))
	}

	if r.r2 != "" {
		res = append(res, fmt.Sprintf("%s  Event Data 2 : 0x%02X, %s ... %s", in, r.ed2, dumpBinary(r.ed2), r.r2))
	} else {
		res = append(res, fmt.Sprintf("%s  Event Data 2 : 0x%02X, %s", in, r.ed2, dumpBinary(r.ed2)))
	}

	if r.r3 != "" {
		res = append(res, fmt.Sprintf("%s  Event Data 3 : 0x%02X, %s ... %s", in, r.ed3, dumpBinary(r.ed3), r.r3))
	} else {
		res = append(res, fmt.Sprintf("%s  Event Data 3 : 0x%02X, %s", in, r.ed3, dumpBinary(r.ed3)))
	}

	if r.genName != "" {
		res = append(res, fmt.Sprintf("%s  Generator    : 0x%04X          ... %s", in, r.genID, r.genName))
	} else {
		res = append(res, fmt.Sprintf("%s  Generator    : 0x%04X", in, r.genID))
	}

	return strings.Join(res, "\n")
}
