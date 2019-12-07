package main

import (
	"fmt"
	"os"

	"github.com/iulianclita/logy/parser"
	"github.com/spf13/cobra"
)

func main() {
	// Flag placeholders
	var (
		text      string
		filter    string
		lines     int
		page      int
		ext       string
		withRegex bool
		noColor   bool
	)
	// Define command
	appCmd := &cobra.Command{
		Use:   "logy /path/to/file",
		Short: "Filter and handle log files of any size with ease",
		Long:  `Filter and handle log files of any size with ease`,
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				if err := cmd.Help(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return
			}
			path := args[0]
			// Create parser object
			p := parser.New(path, text, filter, lines, page, noColor, withRegex, ext)
			// Start parsing the given file
			p.Parse()
		},
	}

	// Parse flags
	appCmd.PersistentFlags().StringVarP(&text, "text", "t", "plain", "Text type to parse (plain/json)")
	appCmd.PersistentFlags().StringVarP(&filter, "filter", "f", "", "Text to filter by")
	appCmd.PersistentFlags().IntVarP(&lines, "lines", "l", 50, "Number of lines per page")
	appCmd.PersistentFlags().IntVarP(&page, "page", "p", 1, "Current page number")
	appCmd.PersistentFlags().StringVarP(&ext, "ext", "e", "", "Accepted file extensions to search in folder")
	appCmd.PersistentFlags().BoolVar(&withRegex, "with-regex", false, "Enable regex support")
	appCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable color output")

	// Run command
	if err := appCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
