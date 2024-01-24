package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
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

func Wget() error {
	// Checking cmd args, if there are none - print advise
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Usage:/wget (url)")
		return fmt.Errorf("not enough arguments")
	}
	// Reading url from cmd args & forming filename with it
	url := args[1]
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

func main() {
	Wget()
}
