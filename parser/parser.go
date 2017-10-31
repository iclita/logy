package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// Parser type definition
type Parser struct {
	path   string
	text   string
	filter string
	lines  int
	page   int
	regex  *regexp.Regexp
	exts   []string
}

// stats for parsed files
type stats struct {
	path    string
	offsets []int64
	matches int
}

// Size for the line scanner buffer
// For very large lines it is safe to put a larger buffer
const scanBuf = 64 * 1024 * 1024

// This is the input format which asks the user for new input data
const inputFmt = "File: %s | Page [%d/%d]\nEnter page number to navigate\nEnter file id and page number separated by a comma to navigate to another file\nPress Ctrl+C if you want to quit:"

const (
	// Graphical display of a checkmark
	yesMark = "\u2714"
	// Graphical display of a multiplication mark
	noMark = "\u2715"
)

// Initial json Regexp
var jsonReg = regexp.MustCompile(`\[?{".*:.*}\]?`)

// Current accepted text types
// This information is useful to let the parser know
// how to interpret the given file input
var textTypes = []string{
	"plain",
	"json",
}

// Color functions to help display meaningful input
var (
	success = color.New(color.FgHiGreen, color.Bold).SprintFunc()
	fail    = color.New(color.FgHiWhite, color.BgRed, color.Bold).SprintFunc()
	alert   = color.New(color.FgHiYellow, color.Bold).SprintFunc()
	info    = color.New(color.FgHiMagenta, color.Bold).SprintFunc()
)

// New returns a new parser object
func New(path, text, filter string, lines, page int, noColor, withRegex bool, ext string) *Parser {
	// Check if a path was provided
	if path == "" {
		exitWithError("Error! Path is required")
	}
	// Check if a valid lines value was provided
	if lines <= 0 {
		exitWithError("Error! Option flag -lines must be strictly positive")
	}
	// Check if a valid page value was provided
	if page <= 0 {
		exitWithError("Error! Option flag -page must be strictly positive")
	}
	// Check if a valid test type was provided
	if !stringInSlice(text, textTypes) {
		exitWithError(fmt.Sprintf("Error! Accepted text types are: %s", strings.Join(textTypes, ", ")))
	}
	// If the path starts with "~"
	// it means the user is probably on a Unix based OS
	// and "~" represents the home directory path
	if strings.HasPrefix(path, "~") {
		user, err := user.Current()
		if err != nil {
			log.Fatal("Cannot get current user:", err)
		}
		path = strings.Replace(path, "~", user.HomeDir, 1)
	}
	// Check if ext flag is present for directory path
	info, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Cannot get file stat info for %s. Error: %v", path, err)
	}
	if info.IsDir() && ext == "" {
		exitWithError("Extensions flag is required for directory paths")
	}
	// Disables colorized output
	if noColor {
		color.NoColor = true
	}
	// Enable regex support is user asks for it
	var regex *regexp.Regexp
	// Attach regexp to parser
	if withRegex {
		// If no filter is provided then regex is useless
		// In this case notify the user
		if filter == "" {
			exitWithError("Error! Regex is present but no filter value was provided!")
		}
		// Enable regex support for filter of length greater than 1
		// If regex remains enabled when filter lenght is 1, strange output is given
		// Also there is no sense in having a regex with length of 1
		if len(filter) > 1 {
			// Compile regex expression here to be user later in the parser
			re, err := regexp.Compile(filter)
			if err != nil {
				exitWithError(fmt.Sprintf("Regex parse error: %s", err.Error()))
			}
			regex = re
		}
	}
	// Get slice with all extensions
	exts := strings.Split(ext, ",")

	return &Parser{
		path:   path,
		text:   text,
		filter: filter,
		lines:  lines,
		page:   page,
		regex:  regex,
		exts:   exts,
	}
}

// Parse parses the file and shows the output to the user
func (p *Parser) Parse() {
	// Channel to receive all file stats
	sc := make(chan stats)
	// File stats slice
	var fs []stats
	// Get all file paths to traverse
	paths := p.getPaths()
	// Get number of possible paths
	numPaths := len(paths)
	// This should never happen :)
	if numPaths == 0 {
		exitWithError("No valid paths were found!")
	}
	// Set current page to what the parser gave us
	currentPage := p.page
	// Count all the lines and provide all page offsets
	// These offsets help us navigate to any page instantly
	// We will do this concurrently :)
	// This is where all the "magic" happens
	for _, path := range paths {
		go p.countLines(path, sc)
	}
	// Wait to receive all file offsets
	for i := 0; i < numPaths; i++ {
		stat := <-sc
		if len(stat.offsets) > 0 {
			fs = append(fs, stat)
		}
	}
	// Readjust with files that have more than 0 pages
	numPaths = len(fs)
	// If nothing can be shown tell the user
	if numPaths == 0 {
		fmt.Printf("%s\n", info("Sorry. Nothing to show here!"))
		return
	}
	// Start from the first ID
	currentID := 1
	// Get first file stat and define it as current stat
	currentFile := fs[currentID-1]
	// Define current path
	currentPath := currentFile.path
	// Define starting offsets for the first file
	currentOffsets := currentFile.offsets
	// Determine total number of pages
	numPages := len(currentOffsets)
	// If nothing can be shown tell the user
	if numPages == 0 {
		fmt.Printf("%s\n", info("Sorry. Nothing to show here!"))
		return
	}
	// The current page cannot be greater than the total number of pages
	if currentPage > numPages {
		fmt.Printf("\n%s\n\n", fail(fmt.Sprintf("Error! Page number cannot be greater than %d", numPages)))
		return
	}
	// Render the table with file stats
	renderStats(fs, currentID)
	fmt.Println()
	// Get the first page output and print it
	output := p.getFilePage(currentPath, currentOffsets[currentPage-1])
	// Send output to the console
	fmt.Println(output)
	// If we haf only 1 file and no more pages are to be shown stop here
	// If means we only have 1 page which we already displayed
	if numPaths == 1 && numPages == 1 {
		return
	}
	// Show a message telling the user at which page we are right now and prompt to navigate to whatever page
	fmt.Print(alert(fmt.Sprintf(inputFmt, currentPath, currentPage, numPages)), " ")
	// Start a new input scanner
	in := bufio.NewScanner(os.Stdin)
	// Scan for incoming input
	for in.Scan() {
		// We only accept page numbers here
		// This is how the parser knows where to navigate next
		id, page, err := extractNavigation(in.Text())
		if err != nil {
			fmt.Printf("\n%s\n\n", fail("Error! A valid number is required"))
			fmt.Print(alert(fmt.Sprintf(inputFmt, currentPath, currentPage, numPages)), " ")
		} else {
			switch {
			case page < 1 || page > numPages:
				fmt.Printf("\n%s\n\n", fail(fmt.Sprintf("Error! Page number must be between 1 and %d", numPages)))
				fmt.Print(alert(fmt.Sprintf(inputFmt, currentPath, currentPage, numPages)), " ")
			case id > numPaths:
				fmt.Printf("\n%s\n\n", fail(fmt.Sprintf("Error! ID number must be between 1 and %d", numPaths)))
				fmt.Print(alert(fmt.Sprintf(inputFmt, currentPath, currentPage, numPages)), " ")
			default:
				// Only if id > 0 then the user wants to change the current file
				// If the only wants to change the page it sends id=0
				if id > 0 {
					currentID = id
					// Set current file stat value
					currentFile = fs[currentID-1]
					// Set current path
					currentPath = currentFile.path
					// Set offsets for the first file
					currentOffsets = currentFile.offsets
					// Determine total number of pages
					numPages = len(currentOffsets)
				}
				// Set current page
				currentPage = page
				fmt.Println()
				// Render table with files
				renderStats(fs, currentID)
				fmt.Println()
				// Get output for display
				output := p.getFilePage(currentPath, currentOffsets[currentPage-1])
				fmt.Printf("\n%s\n", output)
				fmt.Print(alert(fmt.Sprintf(inputFmt, currentPath, currentPage, numPages)), " ")
			}
		}
	}

	if err := in.Err(); err != nil {
		log.Fatal("Scanner error:", err)
	}
}

// getFilePage gets the output for a new page on the input file
func (p *Parser) getFilePage(path string, offset int64) string {
	// Open the file
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open file path %s, Error: %v", path, err)
	}
	defer f.Close()
	// Navigate to the given offset
	// This way we skip the part we don't need
	// and avoid parsing unnecessary lines
	_, err = f.Seek(offset, io.SeekStart)
	if err != nil {
		log.Fatal("Cannot open file stats:", err)
	}
	// Start a new scanner
	s := bufio.NewScanner(f)
	// Set a larger buffer just in case
	s.Buffer(nil, scanBuf)
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
		log.Fatal("File page scanner error", err)
	}

	return output.String()
}

// countLines extracts all line offsets
// This is where all the "magic" happens
// Here we count all the file lines
// and extract a slice with all page offsets
func (p *Parser) countLines(path string, ch chan stats) {
	// Open the file
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open file path %s, Error: %v", path, err)
	}
	defer f.Close()
	// Get file size the know the position
	// in the file of the last byte
	fi, err := f.Stat()
	if err != nil {
		log.Fatal("Cannot open file stats:", err)
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
	// This is used to know if page has been hit (matched)
	// by a filter provided by the user
	var pageHit bool
	// This is used to count the number of matches
	var matches int
	// We start by adding the first page offset which is 0
	pageOffsets = append(pageOffsets, offset)
	// Read all lines one by one
	for {
		currentLine++
		// Read the input file line by line
		// This may take a while depending on the file size
		// This is the most time consuming portion of the app
		line, err := r.ReadBytes('\n')
		// If we have at least 1 line hit it means we have a page hit
		// We also keep track of the total number of hits
		numHits := p.lineHits(line)
		if numHits > 0 {
			pageHit = true
			matches += numHits
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
			// Reset all counters and also the page hit
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
			// This will hold the offsets to be sent on the channel
			var finalOffsets []int64
			// If the input was filtered return filter page offsets
			// Otherwise return normal page offsets
			if p.filter != "" {
				finalOffsets = filterOffsets
			} else {
				finalOffsets = pageOffsets
			}
			ch <- stats{
				path:    path,
				offsets: finalOffsets,
				matches: matches,
			}
			// Exit the function to avoid goroutine leak
			return

		case err != nil:
			log.Fatal("Count lines error:", err)
		}
	}
}

// lineHits determines the number of line matches for a given filter
func (p *Parser) lineHits(line []byte) int {
	// If no filter was provided then we do not care about this
	if p.filter == "" {
		return 0
	}
	// If regex was enbled, search by regex
	// Otherwise search normally by text bytes
	if p.regex != nil {
		matches := p.regex.FindAll(line, -1)
		return len(matches)
	}

	return bytes.Count(line, []byte(p.filter))
}

// getOutput computes the final output
func (p *Parser) getOutput(text string) string {
	// Format input as JSON if needed
	if p.text == "json" {
		jsonMatches := jsonReg.FindAllString(text, -1)
		if jsonMatches != nil {
			for _, m := range jsonMatches {
				text = strings.Replace(text, m, formatJSON(m), -1)
			}
		}
	}
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

// getPaths retrieves file paths for a given root
func (p *Parser) getPaths() []string {
	// Define the final paths
	var paths []string
	// Walk the file/directory
	err := filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("File walk error for: %s", path)
		}
		// If the root path is a file
		// just append use the path without verifying extentions
		if p.path == path && !info.IsDir() {
			paths = append(paths, path)
		} else {
			if info.IsDir() {
				return nil
			}
			ext := strings.Replace(filepath.Ext(path), ".", "", 1)
			// Check if current extension matches one of the desired extensions
			if stringInSlice(ext, p.exts) {
				paths = append(paths, path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal("Directory walk error: ", err)
	}

	return paths
}

// extractNavigation fetches navigation details (file id and page number)
func extractNavigation(s string) (int, int, error) {
	nums := strings.Split(s, ",")
	if len(nums) > 2 {
		exitWithError("More than 2 numbers provided")
	}
	// Only page number counts in this situation
	// The user wants to use the current file and only change the page
	// So we return id=0 as a convention here
	if len(nums) == 1 {
		page, err := strconv.Atoi(strings.Trim(nums[0], " "))
		if err != nil {
			return 0, 0, err
		}
		return 0, page, nil
	}
	// The user wants to change the file and page implicitly
	id, err := strconv.Atoi(strings.Trim(nums[0], " "))
	if err != nil {
		return 0, 0, err
	}
	page, err := strconv.Atoi(strings.Trim(nums[1], " "))
	if err != nil {
		return 0, 0, err
	}
	return id, page, nil
}

// renderStats Displays the current stats for all files
func renderStats(fs []stats, id int) {
	// Set table options
	table := tablewriter.NewWriter(os.Stdout)

	table.SetRowLine(true)
	table.SetCenterSeparator("+")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")
	// Define table header
	table.SetHeader([]string{"File ID", "File Path", "Number of Pages", "Number of Matches", "Current"})
	// Compute the table
	var current string
	for k, v := range fs {
		if id == k+1 {
			current = yesMark
		} else {
			current = noMark
		}
		row := []string{
			strconv.Itoa(k + 1),
			v.path,
			strconv.Itoa(len(v.offsets)),
			strconv.Itoa(v.matches),
			current,
		}
		table.Append(row)
	}
	fmt.Println(info(fmt.Sprintf("Current File ID is %d", id)))
	// Render the table
	table.Render()
}
