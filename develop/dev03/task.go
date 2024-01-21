package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	SortFile("/home/vboxuser/go/src/github.com/Karanth1r3/l_2/develop/dev03/data.txt")
}

func dev03(filePath string) {
	data, err := os.ReadFile("/home/vboxuser/go/src/github.com/Karanth1r3/l_2/develop/dev03/data.txt")
	//	file, err := os.Open(filePath)
	checkError(err)

	fmt.Println(data)
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
		const div = '\n'               // Division char
		line, err := r.ReadString(div) // Reading until meeting div-char
		if err == nil || len(line) > 0 {
			if err != nil {
				line += string(div)
			}
			lines = append(lines, line)
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

func divideByColumns(s string) (result []string) {
	result = strings.Split(s, " ")
	return result
}

func SortFile(filePath string) {

	// Declaring & parsing flags
	var column int
	var path string
	params := os.Args

	if len(params) == 1 {
		fmt.Printf("Usage: ./sort (path_to_file) (flags)")
		return
	}
	paths := strings.Split(filePath, "")
	for _, elem := range path {
		path += 
	}

	reverse := flag.Bool("r", false, "sort in reverse order")
	unique := flag.Bool("u", false, "do not write repeating strings")
	flag.IntVar(&column, "k", 0, "column to sort. columns are divided with space by default")
	flag.Parse()

	filePath = path
	// If u flag is true - only unique strings are going to be read
	lines, err := readLines(filePath)
	if *unique {
		lines = removeDuplicates(lines)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// If reverse order flag is true - sort in reverse order
	if *reverse {
		sort.Sort(sort.Reverse(sort.StringSlice(lines)))
	} else {
		sort.Strings(lines)
	}

	fmt.Println(lines)
	err = writeLines(filePath, lines)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
