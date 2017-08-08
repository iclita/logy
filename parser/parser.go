package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	textTypes = []string{
		"plain",
		"json",
		"html",
		"xml",
	}

	currentPage int

	success = color.New(color.FgHiGreen, color.Bold).SprintFunc()

	fail = color.New(color.FgHiWhite, color.BgRed).SprintFunc()

	alert = color.New(color.FgHiBlack, color.BgHiYellow, color.Bold).SprintFunc()

	info = color.New(color.FgHiBlack, color.BgHiCyan, color.Bold).SprintFunc()
)

type parser struct {
	file    string
	text    string
	filter  string
	lines   int
	page    int
	regex   bool
	noColor bool
}

// New returns a new parser
func New(file, text, filter string, lines, page int, regex, noColor bool) *parser {
	return &parser{
		file:    file,
		text:    text,
		filter:  filter,
		lines:   lines,
		page:    page,
		regex:   regex,
		noColor: noColor,
	}
}

// Parse parses the file and shows the output to the user
func (p *parser) Parse() {
	if p.noColor {
		color.NoColor = true // disables colorized output
	}

	if !stringInSlice(p.text, textTypes) {
		exitWithError(fail(fmt.Sprintf("Error! Accepted text types are: %s", strings.Join(textTypes, ", "))))
	}

	if p.file == "" {
		exitWithError(fail("Error! File is required"))
	}

	currentPage := p.page

	pageOffsets, err := p.countLines()

	if err != nil {
		log.Fatal(err)
	}

	numPages := len(pageOffsets)

	if numPages == 0 {
		fmt.Println(info("Sorry. Nothing to show here!"))
		return
	}

	output := p.getFilePage(pageOffsets[currentPage-1])

	fmt.Println(output)

	if numPages <= 1 {
		return
	}

	fmt.Print(alert(fmt.Sprintf("Page %d/%d. Enter page number to navigate or press Ctrl+C to quit:", currentPage, numPages)), " ")

	in := bufio.NewScanner(os.Stdin)

	for in.Scan() {
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

func (p *parser) getFilePage(offset int64) string {

	f, err := os.Open(p.file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.Seek(offset, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	s := bufio.NewScanner(f)

	var current int
	var output bytes.Buffer
	var text string

	for s.Scan() {
		current++
		if p.filter != "" {
			text = strings.Replace(s.Text(), p.filter, success(p.filter), -1)
		} else {
			text = s.Text()
		}
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

func (p *parser) countLines() ([]int64, error) {

	f, err := os.Open(p.file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	lastPosition := fi.Size()

	r := bufio.NewReader(f)

	var pageOffsets []int64
	var filterOffsets []int64
	var currentLine int
	var lastOffset int64
	var offset int64
	var search []byte
	var pageHit bool

	pageOffsets = append(pageOffsets, offset)

	if p.filter != "" {
		search = []byte(p.filter)
	}

	for {
		currentLine++
		bs, err := r.ReadBytes('\n')
		if search != nil && bytes.Contains(bs, search) {
			pageHit = true
		}
		offset += int64(len(bs))
		if currentLine == p.lines {
			lastOffset = pageOffsets[len(pageOffsets)-1]
			if pageHit {
				filterOffsets = append(filterOffsets, pageOffsets[len(pageOffsets)-1])
			}
			offset += lastOffset
			if offset < lastPosition {
				pageOffsets = append(pageOffsets, offset)
			}
			currentLine = 0
			offset = 0
			pageHit = false
		}

		switch {
		case err == io.EOF:
			if currentLine < p.lines && pageHit {
				filterOffsets = append(filterOffsets, pageOffsets[len(pageOffsets)-1])
			}
			if search != nil {
				return filterOffsets, nil
			}
			return pageOffsets, nil

		case err != nil:
			return nil, err
		}
	}
}

func stringInSlice(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func exitWithError(s string) {
	io.WriteString(os.Stderr, s)
	os.Exit(1)
}
