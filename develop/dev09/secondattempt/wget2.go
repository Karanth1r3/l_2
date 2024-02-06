package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type (
	Wgetter struct {
		links    []string
		filePath string
	}

	WGetOpts struct {
		fileName string
	}
)

func main() {
	ParseFlags(os.Args)
}

func (w *Wgetter) ParseFlags(call []string) error {

	var link string

	if len(call) == 1 {
		fmt.Println("Usage: /wget (flags) link")
		return fmt.Errorf("not enough arguments")
	}
	for _, elem := range call {
		if elem[0] != '-' {
			link = elem
		}
	}
	path := flag.String("f", "", "Path of the saved file")

	flag.Parse()
	if *path != "" {
		w.filePath = *path
	}
	return nil
}

func WGetCli(call []string) error {
	//Piping
	inPipe := os.Stdin
	outPipe := os.Stdout
	errPipe := os.Stderr
	var wget Wgetter
	err := wget.ParseFlags(os.Args)
	if err != nil {
		return err
	}
	return wget.Exec(inPipe, outPipe, errPipe)
}

func (w *Wgetter) Exec(inPipe io.Reader, outPipe io.Writer, errPipe io.Writer) error {
	if len(w.links) > 0 {
		for _, link := range w.links {
			err := WGetOne(link, nil)
			if err != nil {
				return err, 1
			}
		}
	} else {
		// try to read from stdin
		bio := bufio.NewReader(inPipe)
		hasMoreLine := true
		var err error
		var line []byte
		for hasMoreLine {
			line, hasMoreLine, err = bio.ReadLine()
			if err == nil {
				err = WGetOne(strings.TrimSpace(string(line)), nil)
				if err != nil {
					return err, 1
				}
			} else {
				hasMoreLine = false
			}
		}
	}

	return nil, 0
}

func (w *Wgetter) WGetOne(url string) error {

	/*
		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}
	*/
	//resp, err :=

	return nil
}
