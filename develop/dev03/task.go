package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
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

func readLines(filePath string) (lines []string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		const div = '\n'
		line, err := r.ReadString(div)
		if err == nil || len(line) > 0 {
			if err != nil {
				line += string(div)
			}
			lines = append(lines, line)
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
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
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

func SortFile(filePath string) {
	filePath = "/home/vboxuser/go/src/github.com/Karanth1r3/l_2/develop/dev03/data.txt"
	lines, err := readLines(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sort.Strings(lines)
	err = writeLines(filePath, lines)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
