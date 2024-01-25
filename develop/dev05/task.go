package dev05

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	Grep()
}

func readLines(filePath string) (lines []string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		const div = '\n'             // Division char
		line, _, err := r.ReadLine() // Reading until meeting div-char
		if err == nil || len(line) > 0 {
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}

func writeLines(filePath string, lines []string) (err error) {
	// Opening file, handling error, deferring close
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Initializing writer, iterating through string slice & writing it's content to the file
	w := bufio.NewWriter(f)
	defer w.Flush()
	for _, line := range lines {
		_, err := w.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

// PrintGrep - output func for grep (input data is processed && printed here)
func PrintGrep(lines []string, input string, after, before, context int, ignoreCase, invert, fixed, lineNum, count bool) map[int]string {

	m := make(map[int]string)
	var matched bool
	linesCount := 0
	if context > 0 {
		after = context
		before = context
	}
	var expr string
	if ignoreCase {
		expr = "(?i)" + input
	} else {
		expr = input
	}
	r, err := regexp.Compile(expr) // case insensitive by default with checks in deeper block
	/*
		var baseCondition bool
		var noCompile bool
	*/
	if err != nil {
		fmt.Println(err)
		//	noCompile = true
	}
	for index, line := range lines {
		matched = false
		// If searching for exact matches
		if r.MatchString(line) {
			// No time to think, let it be shit
			if fixed {

				if ignoreCase {
					if !strings.EqualFold(input, line) {
						continue
					} else {
						matched = true
					}
				} else {
					if line != input {
						continue
					} else {
						matched = true
					}
				}
			} else {
				matched = true
			}

		}
		if invert && !matched {
			linesCount++
			m[index] = line
		}
		if !invert && matched {
			linesCount++
			m[index] = line
		}
	}
	// If only matches count is required - print it and drop other output
	if count {
		fmt.Println(linesCount)
		return map[int]string{0: strconv.Itoa(linesCount)} // This crap is for testing purposes, sorry
	}
	indexes := make([]int, len(m))
	idx := 0
	for k := range m {
		indexes[idx] = k
		idx++
	}
	sort.IntSlice.Sort(indexes)
	for k := range indexes {
		printAround(lines, indexes[k], before, after, lineNum)
	}
	return m // FOR TESTING
}

func printAround(lines []string, index, before, after int, lineNum bool) {
	// Before/after/context handle
	// Printing before
	for j := index - 1; j > 0 && before > 0; j-- {
		printLine(lines[j], j, lineNum)
		before--
	}
	// Printing main line
	printLine(lines[index], index, lineNum)
	// Printing after
	for j := index + 1; j < len(lines) && after > 0; j++ {
		printLine(lines[j], j, lineNum)
		after--
	}
}

func printLine(line string, index int, lineNum bool) {
	var s string
	if lineNum {
		s = fmt.Sprintf("%d:", index+1) // Indexing of lines in grep starts with 1 as i recall (it may be false)
		s = fmt.Sprintf("%s%s", s, line)
	} else {
		s = line
	}
	fmt.Printf("%s\n", s)
}

// Grep func to call in main & read input
func Grep() {
	// Checking if cmd arguments are present. If they aren't - Print usage advice
	params := os.Args
	if len(params) == 1 {
		fmt.Printf("Usage: %s (flags) (string_pattern) (path_to_file)", filepath.Base(os.Args[0]))
		return
	}

	// Fnal path variable
	var filePath string
	var firstIndex, secondIndex int
	// Name of the file is right after the flags
	for i := 1; i < len(params); i++ {
		if !strings.HasPrefix(params[i], "-") {
			firstIndex = i
			break
		}
	}
	secondIndex = firstIndex + 1
	input := os.Args[firstIndex]
	file := os.Args[secondIndex]
	// Getting current directory
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Splitting & reassembling path for cross-OS support
	pathSplit := filepath.SplitList(path)
	for _, dir := range pathSplit {
		filePath = filepath.Join(dir, file)
	}

	// Declaring & parsing flags
	var before, after, context int
	ignoreCase := flag.Bool("i", false, "ignore case")
	invert := flag.Bool("v", false, "exclude instead of matching")
	fixed := flag.Bool("F", false, "find only exact match to string")
	lineNum := flag.Bool("n", false, "print line number")
	count := flag.Bool("c", false, "print lines count")
	flag.IntVar(&after, "A", 0, "print n lines after match")
	flag.IntVar(&before, "B", 0, "print n lines before match")
	flag.IntVar(&context, "C", 0, "print n lines before and after match")

	//	flag.IntVar(&column, "k", 0, "column to sort. columns are divided with space by default")
	flag.Parse()
	//filePath = path
	// If u flag is true - only unique strings are going to be read
	lines, err := readLines(filePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	PrintGrep(lines, input, after, before, context, *ignoreCase, *invert, *fixed, *lineNum, *count)

}
