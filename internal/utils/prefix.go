package utils

import "strings"

// AddPrefix adds the given prefix to the string based on file format, to save in certain R2 paths.

func AddPrefix(str string) string {
	if strings.HasSuffix(strings.ToLower(str), ".oma") {
		return "omafiles/" + str
	} else {
		return "maps/" + str
	}
}