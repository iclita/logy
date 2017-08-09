package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"bitbucket.org/iulianclita/logy/parser"
	"github.com/olekukonko/tablewriter"
)

// func init() {
// 	b, err := ioutil.ReadFile("test.log")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	f, err := os.OpenFile("plain.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()

// 	for i := 0; i < 100000; i++ {
// 		_, err := f.Write(b)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }

func showHelp() {
	fmt.Println("Welcome to logy. The best parser for filtering and handling log files of any size with ease.")
	fmt.Println("Below is a table explaining the usage of this little utility.")

	data := [][]string{
		[]string{"-file", "Log file path", "logy -file=path/to/file.log", "YES"},
		[]string{"-text", "Text type to parse. Defaults to plain. Valid options are: plain, json, html, xml", "logy -file=path/to/file.log -text=json", "NO"},
		[]string{"-filter", "Text to filter by", "logy -file=path/to/file.log -filter=search", "NO"},
		[]string{"-lines", "Number of lines per page. Defaults to 100", "logy -file=path/to/file.log -lines=250", "NO"},
		[]string{"-page", "Current page number. Defaults to 1", "logy -file=path/to/file.log -page=10", "NO"},
		[]string{"--with-regex", "Enable regex support. Defaults to false", "logy -file=path/to/file.log -filter=[0-9]+search --with-regex", "NO"},
		[]string{"--no-color", "Disable color output. Defaults to false", "logy -file=path/to/file.log --no-color", "NO"},
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
	withRegex := flag.Bool("with-regex", false, "Enable regex support. Defaults to false")
	noColor := flag.Bool("no-color", false, "Disable color output. Defaults to false")

	flag.Parse()

	fmt.Println()

	if len(os.Args) == 1 {
		showHelp()
		return
	}

	var regex *regexp.Regexp

	if *withRegex {
		if *filter == "" {
			log.Fatal("Error! Filter option is missing!")
		}
		re, err := regexp.Compile(*filter)
		if err != nil {
			log.Fatal(err)
		}
		regex = re
	}

	p := parser.New(*file, *text, *filter, *lines, *page, *noColor, regex)

	p.Parse()

}
