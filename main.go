package main

import (
	"log"
	"os"

	"github.com/iulianclita/logy/parser"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "logy",
		Usage:       "smart parsing",
		UsageText:   "logy -lines=10 -filter=error /path/to/file",
		Description: "Filter and handle log files of any size with ease.",
		HideVersion: true,
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
			if c.NArg() == 0 {
				cli.ShowAppHelp(c)
				return nil
			}
			// Extract positional arguments
			path := c.Args().First()
			// Extract flags
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
