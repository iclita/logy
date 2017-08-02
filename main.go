package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var textTypes = []string{
	"plain",
	"json",
	"html",
}

var currentPage int

func lineCounter(file string) (int, error) {

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf := make([]byte, 32*1024)
	count := 1
	lineSep := []byte{'\n'}

	for {
		c, err := f.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func exitWithError(s string) {
	io.WriteString(os.Stderr, s)
	os.Exit(1)
}

func parseFile(file string, page, lines int) string {

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	current := 0

	var output bytes.Buffer

	skip := (page - 1) * lines

	for s.Scan() {
		if current >= skip {
			output.WriteString(s.Text() + "\n")
		}
		current++
		if current >= lines+skip {
			break
		}
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	return output.String()
}

func main() {
	file := flag.String("file", "", "File path")
	text := flag.String("text", "plain", "Text type to parse. Defaults to plain")
	lines := flag.Int("lines", 100, "Number of lines per page. Defaults to 100")
	page := flag.Int("page", 1, "Current page number. Defaults to 1")
	noColor := flag.Bool("no-color", false, "Disable color output. Defaults to false")
	// filter := flag.String("filter", "", "Text to filter")

	flag.Parse()

	if *noColor {
		color.NoColor = true // disables colorized output
	}

	// success := color.New(color.FgGreen, color.Bold).SprintFunc()

	fail := color.New(color.FgWhite, color.BgRed).SprintFunc()

	alert := color.New(color.FgBlack, color.BgYellow).SprintFunc()

	info := color.New(color.FgCyan, color.Bold).SprintFunc()

	if !stringInSlice(*text, textTypes) {
		exitWithError(fail(fmt.Sprintf("Error! Accepted text types are: %s", strings.Join(textTypes, ", "))))
	}

	if *file == "" {
		exitWithError(fail("Error! File is required"))
	}

	currentPage := *page

	n, err := lineCounter(*file)

	if err != nil {
		log.Fatal(err)
	}

	pages := float64(n) / float64(*lines)

	pages = math.Ceil(pages)

	numPages := int(pages)

	fmt.Println(info("Welcome to logy. The best parser for filtering and handling log files with ease.\n"))

	output := parseFile(*file, currentPage, *lines)

	fmt.Println(output)

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
				output := parseFile(*file, currentPage, *lines)
				fmt.Printf("\n%s\n", output)
				fmt.Print(alert(fmt.Sprintf("Page %d/%d. Enter page number to navigate or press Ctrl+C to quit:", currentPage, numPages)), " ")
			}
		}
	}

	if err := in.Err(); err != nil {
		log.Fatal(err)
	}

}
