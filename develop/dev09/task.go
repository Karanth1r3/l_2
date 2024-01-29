package dev09

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func download(url, file string) error {
	// Trying to get response with body
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not get response from url: %w", err)
	}
	//After func done => closing
	defer resp.Body.Close()

	// Creating file to write data
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	// Close file after func done to prevent perfromance etc issures
	defer f.Close()
	// Copying bytes from resp body to file
	_, err = io.Copy(f, resp.Body)
	return nil
}

func Wget(url string) error {

	// Reading url from cmd args & forming filename with it
	file := path.Base(url)
	// If file is not present yet - try to download page though url
	if _, err := os.Stat(file); os.IsNotExist(err) {
		err := download(url, file)
		if err != nil {
			panic(err)
		}
		fmt.Println("Web page saved")
		return nil
	} else {
		return fmt.Errorf("Web page was already saved: %w", err)
	}
}

func ParseArgs() ([]string, error) {
	// Checking cmd args, if there are none - print advise
	args := os.Args
	if len(args) == 1 {
		fmt.Printf("Usage: %s(url)", filepath.Base(os.Args[0]))
		return nil, fmt.Errorf("not enough arguments")
	}
	return args, nil
}

func main() {
	args, err := ParseArgs()
	if err != nil {
		log.Fatal(err)
	}
	url := args[1]
	Wget(url)
}
