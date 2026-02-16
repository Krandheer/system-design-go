package main

import (
	"fmt"
	"sync"
)

// producer generates a sequence of numbers and sends them to a channel.
// It closes the channel when it's done.
func producer(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// worker reads numbers from an input channel, squares them, and sends the
// result to an output channel. Each worker runs in its own goroutine.
// This is the "Fan-out" part of our pipeline.
func worker(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n // Squaring the number as our "work"
		}
	}()
	return out
}

// merger takes a slice of input channels and merges their values into a
// single output channel. This is the "Fan-in" part of our pipeline.
func merger(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// output is a helper function that reads from a single input channel
	// and sends the values to the merged output channel.
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close the output channel once all the input
	// channels are closed and their values have been processed.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	// --- Stage 1: The Producer ---
	// Generates the initial stream of data.
	inputChan := producer(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// --- Stage 2: Fan-out ---
	// We start multiple workers to process the data in parallel.
	// Each worker gets its own output channel.
	numWorkers := 3
	workerChans := make([]<-chan int, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workerChans[i] = worker(inputChan)
	}

	// --- Stage 3: Fan-in ---
	// We merge the results from all the workers into a single channel.
	mergedChan := merger(workerChans...)

	// --- Consume the final results ---
	// We read from the merged channel until it's closed.
	for res := range mergedChan {
		fmt.Println(res)
	}
	
	fmt.Println("Pipeline finished.")
}
