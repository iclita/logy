// Package parser parses the given file input
// It displays a paginated version of the file and optionally
// filters the data according to user preferences
// This aspect makes it suitable for inspecting log files of any size
package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	// Current accepted text types
	// This information is useful to let the parser know
	// how to interpret the ginve file input
	textTypes = []string{
		"plain",
		"json",
		"html",
		"xml",
	}

	// Color functions to help display meaningful input
	success = color.New(color.FgHiBlack, color.BgHiGreen, color.Bold).SprintFunc()
	fail    = color.New(color.FgHiWhite, color.BgRed, color.Bold).SprintFunc()
	alert   = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	info    = color.New(color.FgHiMagenta, color.Bold).SprintFunc()
)

// Parser type definition
type parser struct {
	file    string
	text    string
	filter  string
	lines   int
	page    int
	noColor bool
	regex   *regexp.Regexp
}

// New returns a new parser object
func New(file, text, filter string, lines, page int, noColor bool, regex *regexp.Regexp) *parser {
	return &parser{
		file:    file,
		text:    text,
		filter:  filter,
		lines:   lines,
		page:    page,
		noColor: noColor,
		regex:   regex,
	}
}

// Parse parses the file and shows the output to the user
func (p *parser) Parse() {
	// Disables colorized output
	if p.noColor {
		color.NoColor = true
	}
	// Disables regex support for characters of length 1
	// If regex remains enabled strange out is given
	// Also there is no sense in having a regex with length of 1
	if len(p.filter) == 1 {
		p.regex = nil
	}
	// Check if a valid test type was provided
	if !stringInSlice(p.text, textTypes) {
		exitWithError(fail(fmt.Sprintf("Error! Accepted text types are: %s", strings.Join(textTypes, ", "))))
	}
	// Check if a file path was provided
	if p.file == "" {
		exitWithError(fail("Error! File is required"))
	}
	// Set current page to what the parser gave us
	currentPage := p.page
	// Count all the lines and provide all page offsets
	// These offsets help us navigate to any page instantly
	// This is where all the "magic" happens
	pageOffsets, err := p.countLines()

	if err != nil {
		log.Fatal(err)
	}
	// Determine total number of pages
	numPages := len(pageOffsets)
	// If nothing can be shown tell the user
	if numPages == 0 {
		fmt.Println(info("Sorry. Nothing to show here!"))
		return
	}
	// The current page cannot be greater than the total number of pages
	if currentPage > numPages {
		fmt.Printf("\n%s\n\n", fail(fmt.Sprintf("Error! Page number cannot be greater than %d", numPages)))
		return
	}
	// Get the first page output and print it
	output := p.getFilePage(pageOffsets[currentPage-1])

	fmt.Println(output)
	// If no more pages are to be shown stop here
	// If means we only have 1 page which we already displayed
	if numPages == 1 {
		return
	}
	// Show a message telling the user at which page we are right now and prompt to navigate to whatever page
	fmt.Print(alert(fmt.Sprintf("Page %d/%d. Enter page number to navigate or press Ctrl+C to quit:", currentPage, numPages)), " ")
	// Start a new input scanner
	in := bufio.NewScanner(os.Stdin)
	// Scan for incoming input
	for in.Scan() {
		// We only accept page numbers here
		// This is how the parser knows where to navigate next
		inputPage, err := strconv.Atoi(in.Text())
		if err != nil {
			fmt.Printf("\n%s\n\n", fail("Error! A valid number is required"))
			fmt.Print(alert(fmt.Sprintf("Page %d/%d. Enter page number to navigate or press Ctrl+C to quit:", currentPage, numPages)), " ")
		} else {
			if inputPage > numPages {
				fmt.Printf("\n%s\n\n", fail(fmt.Sprintf("Error! Page number cannot be greater than %d", numPages)))
				fmt.Print(alert(fmt.Sprintf("Page %d/%d. Enter page number to navigate or press Ctrl+C to quit:", currentPage, numPages)), " ")
			} else if inputPage < 1 {
				fmt.Printf("\n%s\n\n", fail("Error! Page number cannot be smaller than 1"))
				fmt.Print(alert(fmt.Sprintf("Page %d/%d. Enter page number to navigate or press Ctrl+C to quit:", currentPage, numPages)), " ")
			} else {
				currentPage = inputPage
				output := p.getFilePage(pageOffsets[currentPage-1])
				fmt.Printf("\n%s\n", output)
				fmt.Print(alert(fmt.Sprintf("Page %d/%d. Enter page number to navigate or press Ctrl+C to quit:", currentPage, numPages)), " ")
			}
		}
	}

	if err := in.Err(); err != nil {
		log.Fatal(err)
	}
}

// Gets the output for a new page on the input file
func (p *parser) getFilePage(offset int64) string {
	// Open the file
	f, err := os.Open(p.file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// Navigate to the given offset
	// This way we skip the part we don't need
	// and avoid parsing unnecesary lines
	_, err = f.Seek(offset, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	// Start a new scanner
	s := bufio.NewScanner(f)
	// Current line. We need to keep track
	// of all current line number to determine
	// when a page ends
	var current int
	// This will hold the final output to be shown to the user
	// It is reponsable to display only 1 page
	var output bytes.Buffer
	// The text that is obtained by the "getOutput" method
	var text string
	// Scan the file and extract all lines
	// Stop when we reach the number of lines per page
	// that the user specified
	for s.Scan() {
		current++
		// Get the line output and added in the buffer
		text = p.getOutput(s.Text())
		fmt.Fprintf(&output, "%s\n", text)
		if current >= p.lines {
			break
		}
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	return output.String()
}

// This is where all the "magic" happens
// Here we count all the file lines
// and extract a slice with all page offsets
func (p *parser) countLines() ([]int64, error) {
	// Open the file
	f, err := os.Open(p.file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// Get file size the know the position
	// in the file of the last byte
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	lastPosition := fi.Size()
	// Start a new reader
	r := bufio.NewReader(f)
	// Here we store all page offsets for all pages
	var pageOffsets []int64
	// Here we store page offsets corresponding to filtered text
	var filterOffsets []int64
	// The first line
	// We need to keep track of every line
	// to know when a page ends
	var currentLine int
	// We store the last offset to help us
	// calculate the next offset in the slice
	var lastOffset int64
	// This is the current offset and is
	// calculated line by line
	var offset int64
	// This is used to know if page has been hit
	// by a filter provided by the user
	var pageHit bool
	// We start by adding the first page offset which is 0
	pageOffsets = append(pageOffsets, offset)

	for {
		currentLine++
		// Read the input file line by line
		// This may take a while depending on the file size
		// This is the most time consuming portion of the app
		line, err := r.ReadBytes('\n')
		// If we have a line hit it means we have a page hit
		if p.lineHit(line) {
			pageHit = true
		}
		// Compute the current offset
		offset += int64(len(line))
		if currentLine == p.lines {
			// If we have reached the end of the page
			// Extract the last page offset
			lastOffset = pageOffsets[len(pageOffsets)-1]
			// If we have a page hit (from filtered input)
			// Add it to the filter offsets
			if pageHit {
				filterOffsets = append(filterOffsets, pageOffsets[len(pageOffsets)-1])
			}
			// Add it to the current offset
			// Append this offset to the end of our slice
			offset += lastOffset
			if offset < lastPosition {
				pageOffsets = append(pageOffsets, offset)
			}
			// Reset all counters
			// and also page hit
			currentLine = 0
			offset = 0
			pageHit = false
		}

		switch {
		case err == io.EOF:
			// If we reached the end of file
			// we must take into account also the last lines
			// which were not caught by the page offsets
			if currentLine < p.lines && pageHit {
				filterOffsets = append(filterOffsets, pageOffsets[len(pageOffsets)-1])
			}
			// If the input was filtered return filter page offsets
			// Otherwise return normal page offsets
			if p.filter != "" {
				return filterOffsets, nil
			}
			return pageOffsets, nil

		case err != nil:
			return nil, err
		}
	}
}

// This determines if line matches a given filter
func (p *parser) lineHit(line []byte) bool {
	// If no filter was provided then we do not care about this
	if p.filter == "" {
		return false
	}
	// If regex was enbled, search by regex
	// Otherwise search normally by text bytes
	if p.regex != nil {
		found := p.regex.Find(line)
		return found != nil
	}

	return bytes.Contains(line, []byte(p.filter))
}

// This computes the final output
func (p *parser) getOutput(text string) string {
	// If no filter was provided give the text as it is
	if p.filter == "" {
		return text
	}
	// For regex/normal search highlight filtered input text
	if p.regex != nil {
		// Find all matches
		matches := p.regex.FindAllString(text, -1)
		if matches == nil {
			return text
		}
		// Highlight every match in given input
		for _, m := range matches {
			text = strings.Replace(text, m, success(m), -1)
		}
		return text
	}

	return strings.Replace(text, p.filter, success(p.filter), -1)
}
