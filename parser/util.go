package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

// Formats input as pretty JSON
func formatJSON(text string) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(text), "", strings.Repeat(" ", 2))
	if err != nil {
		log.Fatal("JSON parse error: ", err)
	}
	return prettyJSON.String()
}

// Exit with a nicely colored error message
func exitWithError(s string) {
	io.WriteString(os.Stderr, fmt.Sprintln(fail(s)))
	os.Exit(1)
}
