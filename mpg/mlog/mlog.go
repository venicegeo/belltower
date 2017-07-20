package mlog

import (
	"io"
	"log"

	"github.com/venicegeo/belltower/mpg"
)

// mlog presents a restricted view of the standard log package

// Verbose controls whether or not Debug and Debugf do anything
var Verbose bool

// log.Lshortfile doesn't work: it will always report the source file is mlog.go,
// because the log pavkage has the call depth hard-coded to 2.
// This will break if the log library ever decides to use bit 10 themselves.
const Lshortfile = 1 << 10
const calldepth = 2

func useSource() bool {
	return Flags()&Lshortfile == Lshortfile
}

func Printf(format string, v ...interface{}) {
	if useSource() {
		format = mutils.SourceFile(calldepth) + " " + format
	}
	log.Printf(format, v...)
}

func Print(v ...interface{}) {
	if useSource() {
		vv := []interface{}{
			mutils.SourceFile(calldepth) + " ",
		}
		for _, x := range v {
			vv = append(vv, x)
		}
		log.Print(vv...)
	} else {
		log.Print(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if Verbose {
		Printf(format, v...)
	}
}

func Debug(v ...interface{}) {
	if Verbose {
		log.Print(v...)
	}
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func SetFlags(flags int) {
	log.SetFlags(flags)
}

func Flags() int {
	return log.Flags()
}
