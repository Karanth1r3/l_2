package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var currentDir string

func main() {
	var err error
	currentDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	input := make(chan string)
	done := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go Read(&wg, input, done)
	go Execute(<-input)
	wg.Wait()
}

func changeDirectory(arg string) {
	newPath, err := formDirectory(arg)
	if err != nil {
		fmt.Println("No such file or directory")
		return
	}
	err = os.Chdir(newPath)
	if err != nil {
		fmt.Println(err)
	}
}

func checkPath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	return nil
}

func formDirectory(arg string) (string, error) {
	var finalPath string
	pathSplit := filepath.SplitList(currentDir)
	for _, dir := range pathSplit {
		finalPath = filepath.Join(dir, arg)
	}
	err := checkPath(finalPath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return finalPath, nil
}

func Read(wg *sync.WaitGroup, input chan<- string, doneCh chan struct{}) {
	for {
		var s string
		_, err := fmt.Scanf("%s", &s)
		if err != nil {
			panic(err)
		}
		if s == "/quit" {
			close(input)
			close(doneCh)
			wg.Done()
			return
		}

		input <- s
	}
}

func Execute(command string) {
	for {
		parts := strings.Split(command, " ")
		switch parts[0] {
		case "cd":
			changeDirectory(parts[1])
		default:
			fmt.Println("Unknown command. Available: cd, pwd, echo <args>, -kill <args>, -ps")
		}
	}
}
