package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
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
	go Execute(input, done)
	wg.Wait()
}

func CheckChdir(path string) {
	err := os.Chdir(path)
	if err != nil {
		fmt.Println(err)
	} else {
		currentDir = path
	}
}

func echo(arg string) {
	arg = strings.ReplaceAll(arg, "\"", "")
	fmt.Println(arg)
}

func changeDirectory(arg string) {
	switch arg {
	// If arg is ',' - no action is required - it is supposed to set directory to current
	case ".":
		return
		// if arg is ~ or empty (empty arg creation is handled on higher level) - go to home directory
	case "~", "":
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		CheckChdir(home)
		// If arg is - => switch to parent directory
	case "-":
		parent := filepath.Dir(currentDir)
		CheckChdir(parent)
	default:
		newPath, err := formDirectory(arg)
		if err != nil {
			fmt.Println("No such file or directory")
			return
		}
		CheckChdir(newPath)
	}
	c, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Current directory set to: ", c)
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

func pwd() {
	fmt.Println(currentDir)
}

func Read(wg *sync.WaitGroup, inputCh chan<- string, doneCh chan struct{}) {
	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		if input == "/quit" {
			close(inputCh)
			close(doneCh)
			wg.Done()
			return
		}

		inputCh <- input
	}
}

func killProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("PID %d returned %v \n", pid, err)
	}

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		err = process.Signal(syscall.Signal(0))
		if err != nil {
			fmt.Printf("Process %d terminated gracefully", pid) // Is that so? Should check != vs ==
			return nil
		}
	}

	log.Printf("Failed to terminate pid %d gracefully, sending SIGKILL \n", pid)
	err = process.Signal(syscall.SIGKILL)
	if err != nil {
		//	log.Fatalf("pid %d returned %v", pid, err)
		return fmt.Errorf("PID %d returned %v", pid, err)
	}
	return nil
}

/*
// https://github.com/mitchellh/go-ps/blob/master/process_unix.go - source of info
func getProcesses() error {

	// That's only a snapshot of a process "list" at the moment of function call
	// Opening the folder with active processes
	d, err := os.Open("/proc")
	if err != nil {
		return err
	}
	defer d.Close()

	results := make([]string, 0, 50)
	// Iterating through process files i guess (by names)
	for {
		names, err := d.Readdirnames(10)
		if err == io.EOF { // If all files are read => break reading loop
			break
		}
		if err != nil {
			return err
		}
		// Processes start with digits, only need to include those files
		for _, name := range names {
			if name[0] < 0 && name[0] > '9' {
				continue
			}

			// Ignore errors as processes may not exist anymore
			pid, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}
		}

	}

}
*/
type Process struct {
	pid  int
	path string
}

func getProcessInfo() {
	matches, err := filepath.Glob("/proc/*/exe") // Processes are stored here i guess
	if err != nil {
		fmt.Println(err)
	}
	// Iterating through active processes
	for _, file := range matches {
		target, _ := os.Readlink(file)

		if len(target) > 0 {
			fmt.Printf("Process path: %+v, Process ID: %s \n", target, filepath.Base(filepath.Dir(file)))
			//fmt.Printf("%+v\n", target)
		}
	}
}

func Execute(input chan string, done chan struct{}) {
	for {
		select {
		//Checking input channel
		case command := <-input:
			// Separating command from args and trash
			parts := strings.Fields(command)
			switch parts[0] {
			// Change directory
			case "cd":
				if len(parts) == 1 {
					changeDirectory("")
				} else {
					changeDirectory(parts[1])
				}
				// Display current directory full path
			case "pwd":
				pwd()
				// Write args from echo to SDTOUT
			case "echo":
				if len(parts) == 1 {
					echo("")
				} else {
					echo(parts[1])
				}
				// Kill process with PID (at least try)
			case "kill":
				if len(parts) == 1 {
					fmt.Println("Enter pid to kill as argument to use kill command")
				} else {
					pid, err := strconv.Atoi(parts[1])
					if err != nil {
						fmt.Println("Could not parse pid from args: ", err)
					}
					killProcess(pid)
				}
			case "ps":
				getProcessInfo()
				// If command could not be recognized => Print advice
			default:
				fmt.Println("Unknown command. Available: cd, pwd, echo <args>, -kill <args>, -ps")
			}
		case <-done:
			return
		}
	}
}
