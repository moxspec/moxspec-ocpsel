package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	var sens, gens, eds, inFile, outFile string
	flag.StringVar(&sens, "s", "", "sensor number")
	flag.StringVar(&gens, "g", "", "generator number")
	flag.StringVar(&eds, "e", "", "event data")
	flag.StringVar(&inFile, "f", "", "path to the file of ipmitool sel elist -v")
	flag.StringVar(&outFile, "o", "", "output path")
	flag.Parse()

	if inFile == "" {
		if sens == "" || gens == "" || eds == "" {
			flag.Usage()
			os.Exit(1)
		}
		res, err := decode(sens, gens, eds)
		if err == nil {
			fmt.Println(res.print())
			os.Exit(0)
		}

		fmt.Println(err)
		os.Exit(1)
	}

	// TODO: ioutil will be deprecated
	bod, err := ioutil.ReadFile(inFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	chunkList := parseBody(bod)
	result := parseChunkList(chunkList)

	if outFile != "" {
		err := ioutil.WriteFile(outFile, []byte(strings.Join(result, "\n")), 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	fmt.Println(strings.Join(result, "\n"))
}

func splitEventData(ed uint64) (byte, byte, byte) {
	ed1 := byte((ed >> 16) & 0xFF)
	ed2 := byte((ed >> 8) & 0xFF)
	ed3 := byte(ed & 0xFF)
	return ed1, ed2, ed3
}

func decode(sens, gens, eds string) (record, error) {
	sen, err := strconv.ParseUint(sens, 16, 32)
	if err != nil {
		return record{}, err
	}

	gen, err := strconv.ParseUint(gens, 16, 32)
	if err != nil {
		return record{}, err
	}

	ed, err := strconv.ParseUint(eds, 16, 32)
	if err != nil {
		return record{}, err
	}

	title, d, err := getDecoder(sen, gen)
	if err != nil {
		return record{}, err
	}

	ed1, ed2, ed3 := splitEventData(ed)
	r1, r2, r3 := d(ed1, ed2, ed3)

	return record{
		title: title,
		ed1:   ed1,
		r1:    r1,
		ed2:   ed2,
		r2:    r2,
		ed3:   ed3,
		r3:    r3,
	}, nil
}
