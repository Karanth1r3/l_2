package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
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
	//var stopChan = make(chan os.Signal, 2)
	//	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, os.Kill)

	wg := sync.WaitGroup{}
	wg.Add(1)
	//	<-stopChan // wait for SIGINT
	go read(&wg, input, done)
	go commandsLoop(input, done)
	wg.Wait()

}

func checkChdir(path string) {
	err := os.Chdir(path)
	if err != nil {
		fmt.Println(err)
	} else {
		currentDir = path
	}
}

func echo(arg string) string {
	arg = strings.ReplaceAll(arg, "\"", "")
	fmt.Println(arg)
	return arg
}

func changeDirectory(arg string) (newDir string, err error) {
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
		checkChdir(home)
		// If arg is - => switch to parent directory
	case "-":
		parent := filepath.Dir(currentDir)
		checkChdir(parent)
	default:
		newPath, err := formDirectory(arg)
		if err != nil {
			return "", fmt.Errorf("could not form directory: %w", err)
		}
		checkChdir(newPath)
	}
	c, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not change directory: %w", err)
	}

	fmt.Println("Current directory set to: ", c)
	return c, nil
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

func read(wg *sync.WaitGroup, inputCh chan<- string, doneCh chan struct{}) {
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

// WORKS ONLY IF BUILD FOR FUCK SAKE........................................
func killProcess(pid int) error {
	process, err := os.FindProcess(pid)
	/*
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("PID %d returned %v \n", pid, err)
		}
	*/
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("pid %d returned %v", pid, err)
	}

	if err := process.Kill(); err != nil {
		fmt.Println(err)
	}
	/*

		err = process.Kill()
		//	err = syscall.Kill(-pid, syscall.SIGKILL)
		if err != nil {
			fmt.Println(err, pid)
			fmt.Println(process.Pid)
			return err
		}
		process.Release()
	*/
	/*
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			fmt.Println(process == nil)
			err = process.Signal(syscall.Signal(0))
			fmt.Println(err)
			if err == nil {
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
	*/
	fmt.Println("Killed", pid)
	return nil
}

// Process is simplified process abstraction with id and path to executable
type Process struct {
	pid  int
	path string
}

func forkExec() (int, error) {
	path, err := os.Executable()
	if err != nil {
		return 0, fmt.Errorf("could not get executable info: %w", err)
	}
	fmt.Println(path)

	pid, err := syscall.ForkExec(path, os.Args, &syscall.ProcAttr{
		Dir:   "/",
		Env:   []string{},
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys:   nil,
	})
	if err != nil {
		return 0, fmt.Errorf("could not execute forkExec: %w", err)
	}
	return pid, nil
}

func getProcessInfo() error {

	matches, err := filepath.Glob("/proc/*/exe") // Processes are stored here i guess
	if err != nil {
		return fmt.Errorf("could not get processes info: %w", err)
	}
	id := "-1"
	// Iterating through active processes
	for _, file := range matches {
		target, _ := os.Readlink(file)
		id = filepath.Base(filepath.Dir(file))
		if len(target) > 0 && id[0] > '0' && id[0] < '9' {
			fmt.Println("--------------------------------------------------------------")
			fmt.Printf("Process path: %+v \n, Process ID: %s \n", target, filepath.Base(filepath.Dir(file)))
			//fmt.Printf("%+v\n", target)
		}
	}
	return nil
}

// Fuck
func checkPipe(input string) (res []string) {
	return strings.Split(input, "|")
}

func commandsLoop(input chan string, done chan struct{}) {
	//r, w := io.Pipe()
	//defer w.Close()
	for {
		select {
		//Checking input channel
		case command := <-input:
			cmds := checkPipe(command)
			for _, cmd := range cmds {
				Execute(cmd)
			}
		case <-done:
			return
		}

	}
}

// Execute tries to execute string command by calling according func
func Execute(cmd string) error {
	/*
		defer w.Close()
		buf := make([]byte, 1024)
		_, _ = r.Read(buf)
		// Execution data from previous command can be used in current
		input := string(buf)
		fmt.Println("INPUT", input)
		execRes := ""
	*/
	parts := strings.Fields(cmd)
	switch parts[0] {
	// Change directory
	case "cd":
		if len(parts) == 1 {
			_, err := changeDirectory("")
			if err != nil {
				return fmt.Errorf("cd command failed: %w", err)
			}
			//	execRes = nDir
		} else {
			_, err := changeDirectory(parts[1])
			if err != nil {
				return fmt.Errorf("cd command failed: %w", err)
			}
			//	execRes = nDir
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
				return fmt.Errorf("could not parse pid from args: %w", err)
			}
			err = killProcess(pid)
			if err != nil {
				return fmt.Errorf("could not kill process: %w", err)
			}
		}
	case "ps":
		err := getProcessInfo()
		if err != nil {
			return fmt.Errorf("could not get process info: %w", err)
		}
	case "forkexec":
		id, err := forkExec()
		if err != nil {
			return fmt.Errorf("forkExec() error: %w", err)
		}
		fmt.Println("child process started, pid: ", id)
		// If command could not be recognized => Print advice
	default:
		fmt.Println("Unknown command. Available: cd, pwd, echo <args>, -kill <args>, -ps")
		return nil
	}

	//	w.Write([]byte(execRes))
	return nil
}

/*
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
*/
