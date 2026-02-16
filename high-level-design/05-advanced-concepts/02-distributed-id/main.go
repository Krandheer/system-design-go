package main

import (
	"fmt"
	"sync"
	"time"
)

// Snowflake Config
const (
	epoch          = int64(1609459200000) // 2021-01-01 00:00:00 UTC (Custom Epoch)
	machineIDBits  = 10                   // 1024 Machines
	sequenceBits   = 12                   // 4096 IDs per ms
	machineIDShift = sequenceBits
	timestampShift = sequenceBits + machineIDBits
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits))
)

type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	machineID     int64
	sequence      int64
}

func NewSnowflake(machineID int64) *Snowflake {
	return &Snowflake{
		machineID: machineID,
	}
}

// NextID generates a unique 64-bit ID
func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	if now == s.lastTimestamp {
		// Same millisecond, increment sequence
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// Sequence overflow, wait for next millisecond
			for now <= s.lastTimestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// New millisecond, reset sequence
		s.sequence = 0
	}

	s.lastTimestamp = now

	// Construct ID:
	// | 41-bit Timestamp | 10-bit MachineID | 12-bit Sequence |
	id := ((now - epoch) << timestampShift) |
		(s.machineID << machineIDShift) |
		s.sequence

	return id
}

func main() {
	// Simulate Machine 1 (e.g., in US-East)
	node1 := NewSnowflake(1)
	
	// Simulate Machine 2 (e.g., in EU-West)
	node2 := NewSnowflake(2)

	fmt.Println("--- Generating IDs on Node 1 ---")
	for i := 0; i < 5; i++ {
		id := node1.NextID()
		fmt.Printf("Node 1 ID: %d (Binary: %b)\n", id, id)
	}

	fmt.Println("\n--- Generating IDs on Node 2 ---")
	for i := 0; i < 5; i++ {
		id := node2.NextID()
		fmt.Printf("Node 2 ID: %d (Binary: %b)\n", id, id)
	}

	fmt.Println("\nNotice: Node 2 IDs are completely different even if generated at the same time.")
	fmt.Println("Also, all IDs are roughly increasing over time.")
}
