package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type (
	Wgetter struct {
		link          string
		internalLinks []string
		filePath      string
		visitedLinks  map[string]bool
	}
)

func main() {
	err := WGetCli(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// TODO - adapt this to recursive download
func crawl(link string, opts Wgetter) error {

	if opts.visitedLinks[link] {
		return fmt.Errorf("already visited %s", link)
	}

	opts.visitedLinks[link] = true
	resp, err := http.Get(link)
	if err != nil {
		return fmt.Errorf("could not get response %w", err)
	}
	defer resp.Body.Close()

	linkCounter := 0
	for _, href := range opts.getLinks(resp.Body) {
		// Only internal links
		if len(href) > 0 && string(href[0]) == "/" && href != link {
			//Skip external links which start with //
			if len(href) > 1 && href[1] == '/' {
				continue
			}
			linkCounter++

			err := crawl(href, opts)
			if err != nil {
				return err
			}
			//time.Sleep(time.Second * 1)
		}
	}
	return nil
}

// TODO - adapt it to wget .Func for getting links recursively (if recursive functionality will be added)
func (w *Wgetter) getLinks(body io.Reader) []string {
	var links []string
	x := html.NewTokenizer(body)
	for {
		tt := x.Next()

		switch tt {
		// End case i guess
		case html.ErrorToken:
			// TODO - prevent dublicates in links
			return links
		case html.StartTagToken, html.EndTagToken:
			token := x.Token()
			//WTF is this
			if "a" == token.Data {
				for _, attr := range token.Attr {
					// If token is link => add it
					if attr.Key == "href" {
						links = append(links, attr.Val)
						w.internalLinks = append(w.internalLinks, attr.Val)
					}
				}
			}
		}
	}
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
	if w.link != "" {
		err := w.WGetOne(w.link)
		if err != nil {
			return err
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
				err = w.WGetOne(strings.TrimSpace(string(line)))
				if err != nil {
					return err
				}
			} else {
				hasMoreLine = false
			}
		}
	}

	return nil
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
