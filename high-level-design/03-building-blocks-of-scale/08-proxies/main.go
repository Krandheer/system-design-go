package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// --- Target Server (The Destination) ---
// This represents "google.com" or any external site.
func startTargetServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// The target sees the request coming from the Proxy, not the original client.
		fmt.Printf("[Target Server] Received request from %s\n", r.RemoteAddr)
		fmt.Fprintf(w, "Hello! I see you are connecting via a proxy.")
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// --- Forward Proxy Server ---
// This acts on behalf of the client.
func startProxyServer() {
	// A real proxy handles the CONNECT method or absolute URLs.
	// For this simulation, we'll use a simple "/fetch" endpoint.
	proxyHandler := func(w http.ResponseWriter, r *http.Request) {
		targetURL := r.URL.Query().Get("url")
		if targetURL == "" {
			http.Error(w, "missing 'url' query param", http.StatusBadRequest)
			return
		}

		fmt.Printf("[Forward Proxy] Client requested: %s. Fetching it on their behalf...\n", targetURL)

		// The Proxy makes the request to the target
		resp, err := http.Get(targetURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// The Proxy returns the target's response to the client
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("[Forward Proxy] Got response from target. Sending to client.\n")
		w.Write(body)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/fetch", proxyHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func main() {
	// 1. Start the Target Server (background)
	go startTargetServer()
	// Give it a moment to start
	time.Sleep(500 * time.Millisecond)

	// 2. Start the Proxy Server (background)
	go startProxyServer()
	time.Sleep(500 * time.Millisecond)

	// 3. Client makes a request
	// The client wants to talk to Target (localhost:8081), but it goes through Proxy (localhost:8080)
	fmt.Println("--- Client: Asking Proxy to fetch Target ---")
	
	proxyURL := "http://localhost:8080/fetch?url=http://localhost:8081"
	resp, err := http.Get(proxyURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("\n[Client] Response received: %s\n", string(body))
}
