package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// Formats input as pretty JSON
func formatJSON(text []byte) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, text, "", "\t")
	if err != nil {
		log.Fatal("JSON parse error: ", err)
	}
	return fmt.Sprintf("\n%s", prettyJSON.String())
}

// Exit with a nicely colored error message
func exitWithError(s string) {
	io.WriteString(os.Stderr, fmt.Sprintln(fail(s)))
	os.Exit(1)
}
