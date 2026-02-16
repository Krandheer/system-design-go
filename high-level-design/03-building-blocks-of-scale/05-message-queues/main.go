package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Message represents a task in our queue.
type Message struct {
	ID      int
	Content string
}

// QueueBroker simulates a message broker (like RabbitMQ or SQS).
type QueueBroker struct {
	queue chan Message
	wg    sync.WaitGroup
}

func NewQueueBroker(bufferSize int) *QueueBroker {
	return &QueueBroker{
		// A buffered channel allows producers to send without waiting for receivers
		// up to the buffer limit.
		queue: make(chan Message, bufferSize),
	}
}

// Publish (Producer) sends a message to the queue.
func (qb *QueueBroker) Publish(msg Message) {
	fmt.Printf("[Producer] Sent Order #%d: %s\n", msg.ID, msg.Content)
	qb.queue <- msg
}

// Subscribe (Consumer) starts a worker to process messages.
func (qb *QueueBroker) Subscribe(workerID int) {
	qb.wg.Add(1)
	go func() {
		defer qb.wg.Done()
		for msg := range qb.queue {
			// Simulate processing time (e.g., packing a box, sending an email)
			processTime := time.Duration(rand.Intn(1000)+500) * time.Millisecond
			fmt.Printf("   [Worker %d] Processing Order #%d (taking %v)...\n", workerID, msg.ID, processTime)
			time.Sleep(processTime)
			fmt.Printf("   [Worker %d] DONE Order #%d\n", workerID, msg.ID)
		}
		fmt.Printf("   [Worker %d] Stopping.\n", workerID)
	}()
}

func main() {
	// 1. Create a Broker with a buffer of 100 messages
	broker := NewQueueBroker(100)

	// 2. Start Consumers (Workers)
	// We'll start 3 workers to handle the load in parallel.
	fmt.Println("--- Starting 3 Fulfillment Workers ---")
	broker.Subscribe(1)
	broker.Subscribe(2)
	broker.Subscribe(3)

	// 3. Simulate Producers
	// Orders come in VERY fast (faster than workers can handle).
	fmt.Println("\n--- Receiving Rush of Orders ---")
	for i := 1; i <= 10; i++ {
		broker.Publish(Message{
			ID:      i,
			Content: fmt.Sprintf("Pack Item SKU-%d", i*100),
		})
		time.Sleep(100 * time.Millisecond) // Fast incoming requests
	}

	// 4. Shutdown
	// We close the queue to signal workers that no more jobs are coming.
	fmt.Println("\n--- No more orders. Waiting for workers to finish... ---")
	close(broker.queue)
	
	// Wait for all workers to finish draining the queue.
	broker.wg.Wait()
	fmt.Println("All orders processed.")
}
