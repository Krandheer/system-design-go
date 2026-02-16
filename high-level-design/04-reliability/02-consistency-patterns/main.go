package main

import (
	"fmt"
	"sync"
	"time"
)

// DataStore simulates a node in our distributed database.
type DataStore struct {
	Name  string
	Value string
	mu    sync.RWMutex
}

func (ds *DataStore) Write(val string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.Value = val
	fmt.Printf("[%s] Written: '%s'\n", ds.Name, val)
}

func (ds *DataStore) Read() string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.Value
}

func main() {
	// 1. Setup Master and Replica
	master := &DataStore{Name: "Master (Region US)", Value: "v1"}
	replica := &DataStore{Name: "Replica (Region EU)", Value: "v1"}

	var wg sync.WaitGroup

	// 2. Perform a Write to Master
	fmt.Println("--- User Updates Profile ---")
	newValue := "v2"
	master.Write(newValue)

	// 3. Simulate Asynchronous Replication (The "Eventual" part)
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Simulate network lag across the Atlantic
		fmt.Println("... Replication in progress (lag) ...")
		time.Sleep(2 * time.Second)
		
		replica.Write(newValue)
		fmt.Println("... Replication Complete ...")
	}()

	// 4. Simulate a Client Reading from the Replica IMMEDIATELY
	// This represents a user in Europe trying to see the profile update.
	fmt.Println("\n--- Read 1: Immediate Read from Replica ---")
	value := replica.Read()
	fmt.Printf("Client sees: '%s' (Is it new? %v)\n", value, value == newValue)

	if value != newValue {
		fmt.Println("Result: STALE DATA. The system is inconsistent momentarily.")
	}

	// 5. Wait for replication to finish
	wg.Wait()

	// 6. Simulate a Client Reading AFTER convergence
	fmt.Println("\n--- Read 2: Read after delay ---")
	value = replica.Read()
	fmt.Printf("Client sees: '%s' (Is it new? %v)\n", value, value == newValue)
	
	if value == newValue {
		fmt.Println("Result: CONSISTENT. The system has reached eventual consistency.")
	}
}
