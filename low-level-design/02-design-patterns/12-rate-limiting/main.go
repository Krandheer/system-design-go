package main

import (
	"fmt"
	"time"
)

func main() {
	const numRequests = 10
	// We want to process at a rate of 2 requests per second.
	// 1 second / 2 requests = 500 milliseconds per request.
	const rate = 500 * time.Millisecond

	// --- The Ticker ---
	// The ticker will send a value on its channel (ticker.C) every 500ms.
	// We use these ticks as permission tokens to process a request.
	ticker := time.NewTicker(rate)
	defer ticker.Stop() // It's good practice to stop the ticker when done.

	// A channel to simulate incoming requests.
	requests := make(chan int, numRequests)

	// --- Simulate a Burst of Incoming Requests ---
	// In a real system, these would be coming from users, other services, etc.
	// We'll send them all at once to our channel.
	go func() {
		for i := 1; i <= numRequests; i++ {
			requests <- i
		}
		close(requests)
	}()
	
	fmt.Printf("Received a burst of %d requests. Processing at a rate of 1 every %v.\n", numRequests, rate)

	// --- The Rate-Limited Processing Loop ---
	// The `for range` on the requests channel will get a request as soon as one is available.
	for req := range requests {
		// This is the crucial line: `<-ticker.C`
		// This line will BLOCK and wait until the ticker sends its next tick.
		// Since the ticker only ticks every 500ms, this loop can only
		// run, at most, once every 500ms. This enforces our rate limit.
		<-ticker.C
		fmt.Printf("Processing request %d at %v\n", req, time.Now().Format("15:04:05.000"))
	}
	
	fmt.Println("All requests processed.")

	// Note for production systems: While time.Ticker is great for simple cases,
	// the `golang.org/x/time/rate` package provides a more powerful and
	// flexible token bucket implementation that can handle bursts and more
	// complex scenarios. It's the standard for production-grade rate limiting in Go.
}
