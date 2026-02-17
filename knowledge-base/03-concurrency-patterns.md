# Concurrency Patterns: The Engine of High Performance

**Target Audience:** Engineers who want to move beyond simple `go func()` and build robust, scalable concurrent systems.

Go's concurrency model (Goroutines + Channels) is its killer feature. This guide covers the patterns used to manage thousands of concurrent tasks efficiently.

---

## 1. Worker Pool
**Problem:** You have 10,000 jobs. Spawning 10,000 goroutines might crash the system or exhaust resources (DB connections, API limits).
**Solution:** A fixed pool of workers that pull jobs from a queue.

```go
package main

import (
    "fmt"
    "time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("Worker %d started job %d\n", id, j)
        time.Sleep(time.Second) // Simulate work
        results <- j * 2
    }
}

func main() {
    const numJobs = 5
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    // Start 3 workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // Send jobs
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs) // Signal no more jobs

    // Collect results
    for a := 1; a <= numJobs; a++ {
        <-results
    }
}
```

---

## 2. Fan-In / Fan-Out
**Problem:** Processing a pipeline of data is slow if done sequentially.
**Solution:**
*   **Fan-Out:** Distribute work to multiple goroutines.
*   **Fan-In:** Collect results from multiple goroutines into a single channel.

```go
package main

import (
    "fmt"
    "sync"
)

// Fan-Out: Multiple workers reading from same channel
func worker(in <-chan int) <-chan int {
    out := make(chan int)
    // if a function returns a channel (out being returned here) and keep writing to it, then it must be goroutine, 
    // because writing to channel is blocking.
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Fan-In: Merge multiple channels into one
func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    output := func(c <-chan int) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }

    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}

// channel is initialised, later in seprate go-routine values being filled in, workers waiting for those values in channel and they will close their output
// channel when they stop listening to in channel, that is when in channel gets closed. merger listen to c1 and c2 out channel of workers, and merge them till 
// they are exhusted, exhaustion get noticed when these channel are close in worker. 
// Finally we start printing merger out channel values and when to stop we know that on basis of when that channel is closed in merger.
func main() {
    in := make(chan int)
    
    // Start 2 workers (Fan-Out)
    c1 := worker(in)
    c2 := worker(in)

    // Merge results (Fan-In)
    out := merge(c1, c2)

    go func() {
        for i := 0; i < 10; i++ {
            in <- i
        }
        close(in)
    }()

    for n := range out {
        fmt.Println(n)
    }
}
```

---

## 3. Rate Limiting
**Problem:** You are sending requests to an external API that allows only 5 requests per second.
**Solution:** Use a `time.Ticker` to throttle execution.

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    requests := make(chan int, 5)
    for i := 1; i <= 5; i++ {
        requests <- i
    }
    close(requests)

    // Limiter: 1 request every 200ms
    limiter := time.Tick(200 * time.Millisecond)

    for req := range requests {
        <-limiter // Block until tick
        fmt.Println("request", req, time.Now())
    }
}
```

---

## 4. Cancellation (Context)
**Problem:** A user cancels a request, or a timeout occurs. You need to stop all related goroutines immediately to save resources.
**Solution:** Use `context.Context`.

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func operation(ctx context.Context) {
    select {
    case <-time.After(2 * time.Second):
        fmt.Println("Operation finished")
    case <-ctx.Done():
        fmt.Println("Operation cancelled:", ctx.Err())
    }
}

func main() {
    // Create a context that times out after 1 second
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel() // Always call cancel to release resources

    // The operation takes 2 seconds, so it will be cancelled
    operation(ctx)
}
```

---

## 5. ErrGroup (Structured Concurrency)
**Problem:** You spawn multiple goroutines. If *one* fails, you want to cancel *all* of them and return the error.
**Solution:** Use `golang.org/x/sync/errgroup`.

```go
package main

import (
    "context"
    "fmt"
    "golang.org/x/sync/errgroup"
)

func main() {
    g, _ := errgroup.WithContext(context.Background())

    // Task 1: Success
    g.Go(func() error {
        return nil
    })

    // Task 2: Failure
    g.Go(func() error {
        return fmt.Errorf("something went wrong")
    })

    // Wait blocks until all functions return
    if err := g.Wait(); err != nil {
        fmt.Println("Group error:", err)
    }
}
```
