package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	SortFile()
}

func removeDuplicates[T comparable](arr []T) []T {
	if len(arr) < 2 {
		return arr
	}
	allKeys := make(map[T]bool)
	list := []T{}
	for _, elem := range arr {
		if _, value := allKeys[elem]; !value {
			allKeys[elem] = true
			list = append(list, elem)
		}
	}
	return list
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
			/*
				if unique {
					hasValue := false
					for _, elem := range lines {
						if elem == line {
							hasValue = true
						}
					}
					if !hasValue {
						lines = append(lines, line)
					}
				} else {
					lines = append(lines, line)
				}
			*/
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
		_, err := w.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func makeComparedSlice(lines []string, index int) []string {
	comparedLines := make([]string, len(lines))
	for idx, line := range lines {
		parts := strings.Split(line, " ")
		comparedLines[idx] = parts[index]
	}
	return comparedLines
}

func hasColumn(lines []string, index int) bool {
	isOk := true
	for _, line := range lines {
		if index > len(strings.Split(line, " ")) {
			isOk = false
		}
	}
	if !isOk {
		fmt.Println("index is too large for this input data. it will be ignored")
	}
	return isOk
}

func sortNumLess(lines []string, index int) []string {

	if !hasColumn(lines, index) {
		index = -1
	}
	if index == -1 {
		sort.SliceStable(lines, func(i, j int) bool {
			a, err := strconv.ParseInt(lines[i], 10, 0)
			if err != nil {
				return false
			}
			b, err := strconv.ParseInt(lines[j], 10, 0)
			if err != nil {
				return false
			}
			return a < b
		})
	} else {
		sort.SliceStable(lines, func(i, j int) bool {
			a, err := strconv.ParseInt(strings.Split(lines[i], " ")[index], 10, 0)
			if err != nil {
				return false
			}
			b, err := strconv.ParseInt(strings.Split(lines[j], " ")[index], 10, 0)
			if err != nil {
				return false
			}
			return a < b
		})
	}

	return lines
}

func sortNumMore(lines []string, index int) []string {
	if !hasColumn(lines, index) {
		index = -1
	}
	if index == -1 {
		sort.SliceStable(lines, func(i, j int) bool {
			a, err := strconv.ParseInt(lines[i], 10, 0)
			if err != nil {
				return false
			}
			b, err := strconv.ParseInt(lines[j], 10, 0)
			if err != nil {
				return false
			}
			fmt.Println(a, b)
			return a > b
		})
	} else {
		sort.SliceStable(lines, func(i, j int) bool {
			a, err := strconv.ParseInt(strings.Split(lines[i], " ")[index], 10, 0)
			if err != nil {
				return false
			}
			b, err := strconv.ParseInt(strings.Split(lines[j], " ")[index], 10, 0)
			if err != nil {
				return false
			}
			fmt.Println(a, b)
			return a > b
		})
	}

	return lines
}

func sortMore(lines []string, index int) []string {
	if !hasColumn(lines, index) {
		index = -1
	}
	if index != -1 {
		sort.SliceStable(lines, func(i, j int) bool {
			return strings.Compare(strings.Split(lines[i], " ")[index], strings.Split(lines[j], " ")[index]) > 0
		})
	} else {
		sort.SliceStable(lines, func(i, j int) bool {
			return strings.Compare(lines[i], lines[j]) > 0
		})
	}

	return lines
}

func sortLess(lines []string, index int) []string {
	if !hasColumn(lines, index) {
		index = -1
	}
	if index != -1 {
		sort.SliceStable(lines, func(i, j int) bool {
			return strings.Compare(strings.Split(lines[i], " ")[index], strings.Split(lines[j], " ")[index]) <= 0
		})
	} else {
		sort.SliceStable(lines, func(i, j int) bool {
			return strings.Compare(lines[i], lines[j]) <= 0
		})
	}

	return lines
}

func removeD(lines []string, index int) []string {
	comparedLines := make([]string, 0)
	if index == -1 {
		comparedLines = lines
	} else {
		comparedLines = makeComparedSlice(lines, index)
	}

	allKeys := make(map[string]bool)
	list := []string{}
	for idx, elem := range comparedLines {
		if _, ok := allKeys[elem]; !ok {
			allKeys[elem] = true
			list = append(list, lines[idx])
		}
	}
	return list
}

func selectColumn(lines []string, index int) []string {

	for lidx, line := range lines {
		parts := strings.Split(line, " ")
		for i, part := range parts {
			//		contains := false
			if i == index {
				lines[lidx] = part
			}
		}
	}
	return lines
}

func SortFile() {
	// Checking if cmd arguments are present. If they aren't - Print usage advice
	params := os.Args
	if len(params) == 1 {
		fmt.Printf("Usage: ./sort (flags) (path_to_file)")
		return
	}
	var filePath string
	var firstIndex int
	// Name of the file is right after the flags
	for i := 1; i < len(params); i++ {
		if !strings.HasPrefix(params[i], "-") {
			firstIndex = i
			break
		}
	}
	if firstIndex == 0 {
		fmt.Printf("Usage: %s (flags) (path_to_file)", filepath.Base(os.Args[0]))
		return
	}
	file := os.Args[firstIndex]
	// Declaring & parsing flags
	var column int
	reverse := flag.Bool("r", false, "sort in reverse order")
	unique := flag.Bool("u", false, "do not write repeating strings")
	num := flag.Bool("n", false, "sort by numerical value")
	flag.IntVar(&column, "k", -1, "column to sort. columns are divided with space by default")
	flag.Parse()
	// in original - indexes start with 1 i guess; this block handles that
	if column < 1 {
		column = -1
	} else {
		column--
	}
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

	lines, err := readLines(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Fucked up but working
	// Numerical check
	if *num {
		if !*reverse {
			sortNumLess(lines, column)
		} else {
			sortNumMore(lines, column)
		}
		// Non-numerical check
	} else {
		// If reverse order flag is true - sort in reverse order
		if !*reverse {
			sortLess(lines, column)
		} else {
			sortMore(lines, column)
		}
	}
	// If u flag is true - only unique strings are going to be exported
	if *unique {
		lines = removeDuplicates(lines)
	}
	for _, elem := range lines {
		fmt.Println(elem)
	}
	err = writeLines(filePath, lines)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
