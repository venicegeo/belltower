package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	// just to load the component library
	_ "github.com/venicegeo/belltower/components"

	"github.com/venicegeo/belltower/engine"
	"github.com/venicegeo/belltower/mpg/mlog"
)

func main() {
	var parseFlag = flag.Bool("parse", false, "if set, parse the input file but do not execute")
	var verboseFlag = flag.Int("verbose", 0, "0=silent, 1=verbose, 2=verboser")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Printf("No .bt file specified\n")
		flag.Usage()
	}
	filename := flag.Arg(0)

	mlog.Verbose = *verboseFlag

	mlog.Printf("-parse:   %t\n", *parseFlag)
	mlog.Printf("-verbose: %d\n", *verboseFlag)
	mlog.Printf("filename: %s\n", filename)

	byts, err := ioutil.ReadFile(filename)
	check(err)

	lines := string(byts)

	graph, err := engine.ParseDSL(lines)
	check(err)

	if *parseFlag {
		mlog.Printf("---\n%s---\n", graph)
		return
	}

	network, err := engine.NewNetwork(graph)
	check(err)

	err = network.Execute(60 * 10)
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(2)
	}
}
