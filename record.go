package main

import (
	"fmt"
	"strings"
)

type record struct {
	title string
	ed1   byte
	r1    string
	ed2   byte
	r2    string
	ed3   byte
	r3    string
}

func (r record) print() string {
	return r.printIndent("")
}

func (r record) printIndent(in string) string {
	var res []string
	res = append(res, fmt.Sprintf("%s%s", in, r.title))
	res = append(res, fmt.Sprintf("%s  ed1 : 0x%02X, %s ... %s", in, r.ed1, dumpBinary(r.ed1), r.r1))
	res = append(res, fmt.Sprintf("%s  ed2 : 0x%02X, %s ... %s", in, r.ed2, dumpBinary(r.ed2), r.r2))
	res = append(res, fmt.Sprintf("%s  ed3 : 0x%02X, %s ... %s", in, r.ed3, dumpBinary(r.ed3), r.r3))
	return strings.Join(res, "\n")
}
