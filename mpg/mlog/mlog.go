/* Copyright 2017, RadiantBlue Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mlog

import (
	"io"
	"log"

	"github.com/venicegeo/belltower/mpg"
)

// mlog presents a restricted view of the standard log package

// Verbose controls whether or not anything is logged
var Verbose int

const (
	VerboseOff  = 0
	VerboseOn   = 1
	VerboseFull = 2
)

// log.Lshortfile doesn't work: it will always report the source file is mlog.go,
// because the log pavkage has the call depth hard-coded to 2.
// This will break if the log library ever decides to use bit 10 themselves.
const Lshortfile = 1 << 10
const calldepth = 2

func useSource() bool {
	return Flags()&Lshortfile == Lshortfile
}

func Printf(format string, v ...interface{}) {
	if Verbose >= VerboseOn {
		if useSource() {
			format = mutils.SourceFile(calldepth) + " " + format
		}
		log.Printf(format, v...)
	}
}

func Print(v ...interface{}) {
	if Verbose >= VerboseOn {
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
}

func Debugf(format string, v ...interface{}) {
	if Verbose >= VerboseFull {
		Printf(format, v...)
	}
}

func Debug(v ...interface{}) {
	if Verbose >= VerboseFull {
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
