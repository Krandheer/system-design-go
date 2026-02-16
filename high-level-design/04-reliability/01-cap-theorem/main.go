package main

import (
	"fmt"
)

type SystemMode string

const (
	CP SystemMode = "CP (Consistency First)"
	AP SystemMode = "AP (Availability First)"
)

type Node struct {
	Name string
	Data string
}

type DistributedSystem struct {
	NodeA            *Node
	NodeB            *Node
	NetworkConnected bool
	Mode             SystemMode
}

func NewDistributedSystem(mode SystemMode) *DistributedSystem {
	return &DistributedSystem{
		NodeA:            &Node{Name: "Node A", Data: "Initial Data"},
		NodeB:            &Node{Name: "Node B", Data: "Initial Data"},
		NetworkConnected: true,
		Mode:             mode,
	}
}

// Write simulates a client writing data to Node A.
// Node A then tries to replicate it to Node B.
func (ds *DistributedSystem) Write(newData string) {
	fmt.Printf("\n--- Client attempting to write '%s' to Node A ---\n", newData)

	if ds.NetworkConnected {
		// Healthy Network: Both modes behave roughly the same (Success)
		ds.NodeA.Data = newData
		ds.NodeB.Data = newData
		fmt.Println("Success: Data written to Node A and replicated to Node B.")
		return
	}

	// --- Network Partition Logic ---
	fmt.Println("ALERT: Network Partition Detected! Connection to Node B is down.")

	if ds.Mode == CP {
		// CP Mode: Must guarantee consistency. If we can't replicate, we refuse the write.
		fmt.Println("CP Decision: Write REJECTED. Cannot guarantee consistency.")
		fmt.Println("Result: System ensures Node A and Node B do not diverge, but write failed (Availability sacrificed).")
	} else {
		// AP Mode: Must guarantee availability. We accept the write locally.
		ds.NodeA.Data = newData
		fmt.Println("AP Decision: Write ACCEPTED on Node A.")
		fmt.Println("Result: Node A has new data, Node B has old data. System is available, but Inconsistent.")
	}
}

func (ds *DistributedSystem) Read() {
	fmt.Printf("Current State -> Node A: [%s] | Node B: [%s]\n", ds.NodeA.Data, ds.NodeB.Data)
}

func main() {
	// 1. Simulate CP Mode (Like a Bank)
	fmt.Println("=== Simulation 1: CP System (e.g., Banking DB) ===")
	cpSystem := NewDistributedSystem(CP)
	cpSystem.Read()
	
	// Network breaks
	cpSystem.NetworkConnected = false
	cpSystem.Write("Transaction $500")
	cpSystem.Read()

	fmt.Println("------------------------------------------------")

	// 2. Simulate AP Mode (Like a Social Media Feed)
	fmt.Println("=== Simulation 2: AP System (e.g., Twitter Feed) ===")
	apSystem := NewDistributedSystem(AP)
	apSystem.Read()

	// Network breaks
	apSystem.NetworkConnected = false
	apSystem.Write("New Tweet: Hello World!")
	apSystem.Read()
}
