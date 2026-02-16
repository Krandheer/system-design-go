package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

func main() {
	// 1. Allow setting the port via command line argument.
	// This enables "Horizontal Scaling" simulation by running multiple instances on different ports.
	port := flag.String("port", "8080", "Port to run the server on")
	name := flag.String("name", "Server-1", "Name of this server instance")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from %s running on port %s\n", *name, *port)
	})

	// 2. Simulate a heavy CPU task.
	// This represents a request that consumes significant resources.
	http.HandleFunc("/heavy", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Simulate work: Check for primes up to a large number
		count := 0
		for i := 2; i < 50000; i++ {
			if isPrime(i) {
				count++
			}
		}

		duration := time.Since(start)
		msg := fmt.Sprintf("[%s] Heavy calculation done! Found %d primes. Took %v\n", *name, count, duration)
		fmt.Print(msg) // Log to server console
		fmt.Fprint(w, msg) // Send to client
	})

	fmt.Printf("Starting %s on port %s...\n", *name, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

// isPrime is a simple (and purposefully inefficient) CPU-bound function.
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

