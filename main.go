package main

import (
	"log"
	"os"

	"github.com/iulianclita/logy/parser"
	"github.com/urfave/cli/v2"
)

// // Displays help information about this tool
// func showHelp() {
// 	// Show nicely colored info
// 	c := color.New(color.FgHiCyan, color.Bold)
// 	c.Println("Welcome to logy. The best parser for filtering and handling log files of any size with ease.")
// 	c.Println("Below is a table explaining the usage of this little utility.")
// 	// Define table content
// 	data := [][]string{
// 		[]string{"-path", "File/directory path", "logy -file=path/to/file.log OR logy -file=path/to/directory", "YES"},
// 		[]string{"-text", "Text type to parse. Defaults to plain. Valid options are: plain, json", "logy -file=path/to/file.log -text=json", "NO"},
// 		[]string{"-filter", "Text to filter by", "logy -file=path/to/file.log -filter=search", "NO"},
// 		[]string{"-lines", "Number of lines per page. Defaults to 50", "logy -file=path/to/file.log -lines=250", "NO"},
// 		[]string{"-page", "Current page number. Defaults to 1", "logy -file=path/to/file.log -page=10", "NO"},
// 		[]string{"--with-regex", "Enable regex support. Defaults to false", "logy -file=path/to/file.log -filter=[0-9]+search --with-regex", "NO"},
// 		[]string{"-ext", "Accepted file extensions to search in folder. Required only for directory paths", "logy -file=path/to/directory -ext=log,txt", "NO"},
// 		[]string{"--no-color", "Disable color output. Defaults to false", "logy -file=path/to/file.log --no-color", "NO"},
// 	}
// 	// Set table options
// 	table := tablewriter.NewWriter(os.Stdout)

// 	table.SetRowLine(true)
// 	table.SetCenterSeparator("+")
// 	table.SetColumnSeparator("|")
// 	table.SetRowSeparator("-")
// 	// Define table header
// 	table.SetHeader([]string{"Option", "Description", "Usage", "Required"})
// 	// Compute the table
// 	for _, v := range data {
// 		table.Append(v)
// 	}
// 	// Render the table
// 	table.Render()
// }

func main() {
	// fmt.Println()
	// // If no command line flags are provided show help information
	// if len(os.Args) == 1 {
	// 	showHelp()
	// 	return
	// }
	// // If -h or --help command line flags are provided show help information
	// if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
	// 	showHelp()
	// 	return
	// }
	// // Capture and parse incoming commnad line flags
	// path := flag.String("path", "", "File/directory path")
	// text := flag.String("text", "plain", "Text type to parse. Defaults to plain")
	// filter := flag.String("filter", "", "Text to filter by")
	// lines := flag.Int("lines", 50, "Number of lines per page. Defaults to 50")
	// page := flag.Int("page", 1, "Current page number. Defaults to 1")
	// withRegex := flag.Bool("with-regex", false, "Enable regex support. Defaults to false")
	// ext := flag.String("ext", "", "Accepted file extensions to search in folder")
	// noColor := flag.Bool("no-color", false, "Disable color output. Defaults to false")

	// flag.Parse()
	// // Create a new parser object
	// p := parser.New(*path, *text, *filter, *lines, *page, *noColor, *withRegex, *ext)
	// // Start parsing the given file
	// p.Parse()

	app := &cli.App{
		// Define flags
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "text",
				Value: "plain",
				Usage: "Text type to parse (plain/json). Defaults to plain",
			},
			&cli.StringFlag{
				Name:  "filter",
				Value: "",
				Usage: "Text to filter by",
			},
			&cli.IntFlag{
				Name:  "lines",
				Value: 50,
				Usage: "Number of lines per page. Defaults to 50",
			},
			&cli.IntFlag{
				Name:  "page",
				Value: 1,
				Usage: "Current page number. Defaults to 1",
			},
			&cli.BoolFlag{
				Name:  "with-regex",
				Value: false,
				Usage: "Enable regex support. Defaults to false",
			},
			&cli.StringFlag{
				Name:  "ext",
				Value: "",
				Usage: "Accepted file extensions to search in folder",
			},
			&cli.BoolFlag{
				Name:  "no-color",
				Value: false,
				Usage: "Disable color output. Defaults to false",
			},
		},
		// Execute command action
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return nil
			}
			// First (and only) arg should be the path
			path := c.Args().Get(0)
			// These should be the flags
			text := c.String("text")
			filter := c.String("filter")
			lines := c.Int("lines")
			page := c.Int("page")
			withRegex := c.Bool("with-regex")
			ext := c.String("ext")
			noColor := c.Bool("no-color")
			// Create a new parser object
			p := parser.New(path, text, filter, lines, page, noColor, withRegex, ext)
			// Start parsing the given file
			p.Parse()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
