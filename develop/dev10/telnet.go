package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/Karanth1r3/l_2/develop/dev10/telnet"
)

const (
	sec  = "s"
	min  = "m"
	hour = "h"
)

func getDur(s string, value int) (time.Duration, error) {

	switch string(s[len(s)-1]) {
	case sec:
		return time.Second * time.Duration(value), nil
	case min:
		return time.Minute * time.Duration(value), nil
	case hour:
		return time.Hour * time.Duration(value), nil
	default:
		return time.Nanosecond, fmt.Errorf("could not parse duration")
	}

}

func main() {

	// Parsing args (lazy way)
	args := os.Args
	if len(args) < 4 {
		fmt.Println("Usage example:go-telnet --timeout=5s host port")
		return
	}
	host := os.Args[2]
	port := os.Args[3]
	t := os.Args[1]
	fmt.Println(t)
	re := regexp.MustCompile("[0-9]+")
	ts := re.FindAllString(t, -1)

	timeVal, err := strconv.Atoi(ts[0])
	if err != nil {
		log.Fatal("could not parse timeout", err)
	}
	timeout, err := getDur(t, timeVal)
	if err != nil {
		log.Fatal(err)
	}

	// Launching client
	tc, err := telnet.NewTelnetClient(host, port, timeout)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to", tc.Conn.RemoteAddr())
	defer tc.Close()
	//tc.SetTimeout(5)
	//msgCh := make(chan string)

	// Channel reacts after os signals are received in client implementation
	endCh := tc.EndCh
	// Launching client read-write loop (goroutine with goroutines)
	tc.LaunchReadWrite()

	// There is goroutine waiting for os signal and then sending data to endCh
	tc.WaitOSKill()
	// For blocking main until signal is received
	<-endCh

	// If data from endChannel received => stop read-write & close connection
	if err := tc.Cancel(); err != nil {
		log.Println(err)
	}
	// Cleanup time for sockets? (i guess)
	time.Sleep(time.Second)
}
