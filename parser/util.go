package parser

import (
	"io"
	"os"
)

// Checks if a slice of strings contains a given string
func stringInSlice(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// Exit with a nicely colored error message
func exitWithError(s string) {
	io.WriteString(os.Stderr, fail(s)+"\n")
	os.Exit(1)
}
