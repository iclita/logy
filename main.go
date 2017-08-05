package main

import (
	"flag"
	"fmt"
	"os"

	"bitbucket.org/iulianclita/logy/parser"
	"github.com/olekukonko/tablewriter"
)

func showHelp() {
	fmt.Println("Welcome to logy. The best parser for filtering and handling log files with ease.")
	fmt.Println("Below is a table explaining the usage of this little utility.")

	data := [][]string{
		[]string{"file", "Log file path", "logy -file=path/to/file.log", "YES"},
		[]string{"text", "Text type to parse. Defaults to plain. Valid options are: plain, json, html, xml", "logy -file=path/to/file.log -text=json", "NO"},
		[]string{"filter", "Text to filter by", "logy -file=path/to/file.log -filter=search", "NO"},
		[]string{"lines", "Number of lines per page. Defaults to 100", "logy -file=path/to/file.log -lines=250", "NO"},
		[]string{"page", "Current page number. Defaults to 1", "logy -file=path/to/file.log -page=10", "NO"},
		[]string{"regex", "Enable regex support. Defaults to false", "logy -file=path/to/file.log -filter=^[0-9]+search$ -regex", "NO"},
		[]string{"no-color", "Disable color output. Defaults to false", "logy -filter=search --no-color", "NO"},
	}

	table := tablewriter.NewWriter(os.Stdout)

	table.SetRowLine(true)
	table.SetCenterSeparator("+")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")

	table.SetHeader([]string{"Option", "Description", "Usage", "Required"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func main() {
	file := flag.String("file", "", "File path")
	text := flag.String("text", "plain", "Text type to parse. Defaults to plain")
	filter := flag.String("filter", "", "Text to filter by")
	lines := flag.Int("lines", 100, "Number of lines per page. Defaults to 100")
	page := flag.Int("page", 1, "Current page number. Defaults to 1")
	regex := flag.Bool("regex", false, "Enable regex support. Defaults to false")
	noColor := flag.Bool("no-color", false, "Disable color output. Defaults to false")

	flag.Parse()

	if len(os.Args) == 1 {
		showHelp()
		return
	}

	p := parser.New(*file, *text, *filter, *lines, *page, *regex, *noColor)

	p.Parse()

}
