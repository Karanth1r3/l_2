package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	// Create "united" channel instance
	orCh := make(chan interface{})
	// Iterating through all input  channels
	for _, c := range channels {
		// Goroutine will be blocked until something is put into at leas one of the channels
		go func(c <-chan interface{}) {
			//If data is received => no need to put it anywhere, goroutine is blocked until receiving c
			<-c
			// If at least 1 channels is done, close united channel
			close(orCh)
		}(c)
	}

	return orCh
}

// Single signal channel example
func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}
