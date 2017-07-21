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
package merr

import (
	"fmt"

	mutils "github.com/venicegeo/belltower/mpg"
)

var UseSourceStamp = true

type MErr struct {
	Message string
	Source  string
}

func New(mssg string) MErr {
	return MErr{
		Message: mssg,
		Source:  mutils.SourceFile(2),
	}
}

func Newf(format string, v ...interface{}) MErr {
	return MErr{
		Message: fmt.Sprintf(format, v...),
		Source:  mutils.SourceFile(2),
	}
}

func (err MErr) Error() string {
	return err.String()
}

func (err MErr) String() string {
	if UseSourceStamp {
		return err.Source + " " + err.Message
	}
	return err.Message
}
