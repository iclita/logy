package main

import (
	"flag"

	"bitbucket.org/iulianclita/logy/parser"
)

func main() {
	file := flag.String("file", "", "File path")
	text := flag.String("text", "plain", "Text type to parse. Defaults to plain")
	filter := flag.String("filter", "", "Text to filter")
	lines := flag.Int("lines", 100, "Number of lines per page. Defaults to 100")
	page := flag.Int("page", 1, "Current page number. Defaults to 1")
	noColor := flag.Bool("no-color", false, "Disable color output. Defaults to false")

	flag.Parse()

	p := parser.New(*file, *text, *filter, *lines, *page, *noColor)

	p.Parse()

}
