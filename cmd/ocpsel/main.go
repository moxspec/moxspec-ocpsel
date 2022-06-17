package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/moxspec/moxspec-ocpsel/ocpsel"
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
		res, err := ocpsel.Decode(sens, gens, eds)
		if err == nil {
			fmt.Println(res.Summary())
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

	chunkList := ocpsel.ParseText(bod)
	result := ocpsel.ParseChunkList(chunkList)

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
