package main

import (
	"fmt"
	"time"
)

// worker is our concurrent worker. It receives work from the jobs channel
// and sends the corresponding result to the results channel.
// Each worker runs in its own goroutine.
func worker(id int, jobs <-chan int, results chan<- string) {
	// The `for range` on a channel will block until a value is sent to it.
	// When the channel is closed and all values have been received, the loop terminates.
	for j := range jobs {
		fmt.Printf("Worker %d: started job %d\n", id, j)
		
		// Simulate doing some work that takes time.
		time.Sleep(time.Second) 
		
		fmt.Printf("Worker %d: finished job %d\n", id, j)
		
		// Send the result of the work to the results channel.
		results <- fmt.Sprintf("Result of job %d from worker %d", j, id)
	}
}

func main() {
	const numJobs = 10
	const numWorkers = 3

	// Create buffered channels. A buffered channel allows a certain number of values
	// to be sent without a corresponding receiver being ready.
	jobs := make(chan int, numJobs)
	results := make(chan string, numJobs)

	fmt.Printf("Starting a pool of %d workers to handle %d jobs.\n", numWorkers, numJobs)

	// --- Start the Worker Pool ---
	// This launches our fixed number of worker goroutines.
	// They will all block, waiting for work to be sent on the jobs channel.
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// --- Send Jobs ---
	// Here, we send all our jobs to the jobs channel.
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	// We close the jobs channel to signal to the workers that there's no more work.
	// The `for range` loop in each worker will exit after processing the remaining jobs.
	close(jobs)

	// --- Collect Results ---
	// Finally, we collect all the results from the work.
	// We expect to receive one result for each job.
	for a := 1; a <= numJobs; a++ {
		// This will block until a result is available on the channel.
		res := <-results
		fmt.Println("Main: received result ->", res)
	}
	
	fmt.Println("All jobs have been processed.")
}
