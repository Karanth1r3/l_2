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
	//fmt.Println(len(ts))
	timeout, err := getDur(t, timeVal)
	if err != nil {
		log.Fatal(err)
	}

	tc, err := telnet.NewTelnetClient(host, port, timeout)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tc.Conn.RemoteAddr())
	defer tc.Close()
	//tc.SetTimeout(5)
	//msgCh := make(chan string)
	endCh := tc.EndCh
	tc.LaunchReadWrite()

	tc.WaitOSKill()
	<-endCh

	if err := tc.Cancel(); err != nil {
		log.Println(err)
	}

	time.Sleep(time.Second)
}
