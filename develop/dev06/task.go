package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	Cut()
}

// ExecuteCut formats input string data by columns depending on the delimeter
func ExecuteCut(lines []string, indexes []int, delim *string, sepOnly bool) {
	result := [][]string{}

	for _, line := range lines {
		// If -s flag is set and line doesn't include delimiter - ignore it
		if sepOnly && !strings.Contains(line, *delim) {
			continue
		}
		part := strings.Split(line, *delim)
		for i := 0; i < len(part); i++ {
			// Adding separated part of the current string to the according column of result
			if len(result) <= i {
				newColumn := []string{part[i]}
				result = append(result, newColumn)
			} else {
				result[i] = append(result[i], part[i])
			}
		}
	}
	var output string

	for i := 0; i < len(result); i++ {
		// If result doesn't contain column with specified index => ignore it
		contains := false
		for j := 0; j < len(indexes); j++ {
			if i == indexes[j] {
				contains = true
			}
		}
		if !contains {
			continue
		}
		// Form output string
		for _, column := range result[i] {

			output += fmt.Sprintf("%s", column)
		}
	}
	fmt.Printf("%s\n", output)
}

// Cut gets data for cut Cut Execution function (kind of bad)
func Cut() {
	var raw string
	for i := 1; i < len(os.Args); i++ {
		raw += os.Args[i] + " "
	}

	args := strings.Fields(raw)
	flagset := flag.NewFlagSet("FlagSet", flag.ContinueOnError)
	delim := "\t"
	f := flagset.String("f", "", "select fields to export")
	d := flagset.String("d", "", "custom delimiter")
	s := flagset.Bool("s", false, "export only strings with delimeter")
	flagset.Parse(args)
	// If d param is present - set custom delimeter
	if *d != "" {
		delim = (*d)
	}

	idxs := getLineIndexes(*f)

	params := os.Args
	if len(params) == 1 {
		fmt.Printf("Usage: %s (flags) (string data)", filepath.Base(os.Args[0]))
		return
	}

	var firstIndex int
	// Name of the file is the first argument
	for i := 1; i < len(params); i++ {
		if !strings.HasPrefix(params[i], "-") {
			firstIndex = i
			break
		}
	}

	input := os.Args[firstIndex]
	lines := strings.Split(input, "\n")

	ExecuteCut(lines, idxs, &delim, *s)
}

func getLineIndexes(f string) (columns []int) {
	res := []string{}
	s := strings.Split(f, ",")
	for _, arg := range s {
		var temp string
		for _, c := range arg {
			if unicode.IsDigit(c) {
				temp += string(c)
			}
			res = append(res, temp)
		}
	}
	for _, elem := range res {
		val, err := strconv.Atoi(elem)
		if err != nil {
			fmt.Println("could not parse vals from -f flag: ", err)
			continue
		}
		columns = append(columns, val)
	}
	return columns
}
