package ocpsel

import (
	"fmt"
	"strings"
)

// record represents a SEL record
type record struct {
	title string
	ed1   byte
	r1    string
	ed2   byte
	r2    string
	ed3   byte
	r3    string
}

// Summary returns record summary
func (r record) Summary() string {
	return r.SummaryIndent("")
}

// SummaryIndent returns record summary with given indent
func (r record) SummaryIndent(in string) string {
	var res []string

	res = append(res, fmt.Sprintf("%sDeciphered Event Data:", in))
	if r.title != "" {
		res = append(res, fmt.Sprintf("%s  %s", in, r.title))
	}

	if r.r1 != "" {
		res = append(res, fmt.Sprintf("%s  ed1 : 0x%02X, %s ... %s", in, r.ed1, dumpBinary(r.ed1), r.r1))
	} else {
		res = append(res, fmt.Sprintf("%s  ed1 : 0x%02X, %s", in, r.ed1, dumpBinary(r.ed1)))
	}

	if r.r2 != "" {
		res = append(res, fmt.Sprintf("%s  ed2 : 0x%02X, %s ... %s", in, r.ed2, dumpBinary(r.ed2), r.r2))
	} else {
		res = append(res, fmt.Sprintf("%s  ed2 : 0x%02X, %s", in, r.ed2, dumpBinary(r.ed2)))
	}

	if r.r3 != "" {
		res = append(res, fmt.Sprintf("%s  ed3 : 0x%02X, %s ... %s", in, r.ed3, dumpBinary(r.ed3), r.r3))
	} else {
		res = append(res, fmt.Sprintf("%s  ed3 : 0x%02X, %s", in, r.ed3, dumpBinary(r.ed3)))
	}

	return strings.Join(res, "\n")
}
