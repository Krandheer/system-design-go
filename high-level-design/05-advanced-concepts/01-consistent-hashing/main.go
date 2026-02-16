package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

// HashRing handles the consistent hashing logic.
type HashRing struct {
	// Sorted list of hash values (keys) on the ring.
	keys []int
	// Map from hash value to physical node name.
	hashMap map[int]string
	// Number of virtual nodes per physical node.
	replicas int
}

func NewHashRing(replicas int) *HashRing {
	return &HashRing{
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

// AddNode adds a new physical node to the ring.
func (h *HashRing) AddNode(nodeName string) {
	for i := 0; i < h.replicas; i++ {
		// Create virtual node key: "NodeA#1", "NodeA#2", etc.
		virtualNodeKey := nodeName + "#" + strconv.Itoa(i)
		hash := int(crc32.ChecksumIEEE([]byte(virtualNodeKey)))
		
		h.keys = append(h.keys, hash)
		h.hashMap[hash] = nodeName
		fmt.Printf("Added Virtual Node: %s -> Hash: %d\n", virtualNodeKey, hash)
	}
	// Keep the keys sorted for binary search.
	sort.Ints(h.keys)
}

// RemoveNode removes a physical node from the ring.
func (h *HashRing) RemoveNode(nodeName string) {
	for i := 0; i < h.replicas; i++ {
		virtualNodeKey := nodeName + "#" + strconv.Itoa(i)
		hash := int(crc32.ChecksumIEEE([]byte(virtualNodeKey)))
		
		delete(h.hashMap, hash)
		
		// Remove from sorted keys list (linear scan for simplicity)
		for j, k := range h.keys {
			if k == hash {
				h.keys = append(h.keys[:j], h.keys[j+1:]...)
				break
			}
		}
	}
	fmt.Printf("Removed Node: %s\n", nodeName)
}

// GetNode finds the closest node clockwise for a given key.
func (h *HashRing) GetNode(key string) string {
	if len(h.keys) == 0 {
		return ""
	}

	hash := int(crc32.ChecksumIEEE([]byte(key)))
	
	// Binary Search: Find the first key on the ring >= hash
	idx := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hash
	})

	// Wrap around: If we reached the end of the slice, go to the start (0).
	if idx == len(h.keys) {
		idx = 0
	}

	return h.hashMap[h.keys[idx]]
}

func main() {
	// Create a ring with 3 virtual nodes per physical node
	ring := NewHashRing(3)

	// Add 3 physical servers
	ring.AddNode("Server-A")
	ring.AddNode("Server-B")
	ring.AddNode("Server-C")

	fmt.Println("\n--- Distributing Keys ---")
	keys := []string{"User1", "User2", "User3", "User4", "User5"}
	for _, key := range keys {
		node := ring.GetNode(key)
		fmt.Printf("Key '%s' mapped to -> %s\n", key, node)
	}

	fmt.Println("\n--- Removing Server-A (Simulating Crash) ---")
	ring.RemoveNode("Server-A")

	fmt.Println("\n--- Re-Distributing Keys ---")
	for _, key := range keys {
		node := ring.GetNode(key)
		fmt.Printf("Key '%s' mapped to -> %s\n", key, node)
	}
	
	fmt.Println("\nNotice that keys mapped to Server-B and Server-C DID NOT MOVE.")
	fmt.Println("Only keys originally on Server-A were re-assigned.")
}
