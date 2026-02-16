package main

import (
	"fmt"
	"time"
)

// OriginServer simulates our main server in New York.
// It has all the data but is "far away" (slow to access).
type OriginServer struct {
	content map[string]string
}

func (s *OriginServer) Fetch(filename string) (string, bool) {
	fmt.Printf("   [Origin Server] Fetching '%s' from disk...\n", filename)
	// Simulate the latency of traveling across the world (e.g., NYC to Sydney)
	time.Sleep(500 * time.Millisecond)
	
	data, ok := s.content[filename]
	return data, ok
}

// EdgeServer simulates a CDN node in Sydney.
// It is "close" to the user but starts empty.
type EdgeServer struct {
	origin *OriginServer
	cache  map[string]string
}

func (s *EdgeServer) GetFile(filename string) string {
	fmt.Printf("[Edge Server] Request for '%s'...\n", filename)

	// 1. Check local cache
	if data, ok := s.cache[filename]; ok {
		fmt.Println("   -> HIT! Serving from Edge cache (Fast).")
		return data
	}

	// 2. If missing, fetch from Origin
	fmt.Println("   -> MISS. Fetching from Origin (Slow)...")
	data, found := s.origin.Fetch(filename)
	if !found {
		return "404 Not Found"
	}

	// 3. Store in local cache for next time
	s.cache[filename] = data
	return data
}

func main() {
	// Setup the world
	origin := &OriginServer{
		content: map[string]string{
			"index.html": "<html>...</html>",
			"logo.png":   "BINARY_IMAGE_DATA",
		},
	}

	edge := &EdgeServer{
		origin: origin,
		cache:  make(map[string]string),
	}

	// --- User in Sydney requests a file ---
	
	fmt.Println("--- Request 1 (First visit) ---")
	start := time.Now()
	edge.GetFile("logo.png")
	fmt.Printf("Total time: %v\n\n", time.Since(start))

	fmt.Println("--- Request 2 (Refresh page) ---")
	start = time.Now()
	edge.GetFile("logo.png")
	fmt.Printf("Total time: %v\n", time.Since(start))
}
