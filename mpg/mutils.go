package mutils

import (
	"fmt"
	"runtime"
	"strings"
)

func SourceFile(depth int) string {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}

	idx := strings.LastIndex(file, "/")
	// if we returned -1 (not found), the index used will be 0 (cute)

	short := file[idx+1 : len(file)]
	return fmt.Sprintf("%s:%d", short, line)
}
