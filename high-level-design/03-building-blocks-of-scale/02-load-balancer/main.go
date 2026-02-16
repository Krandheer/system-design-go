package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

// ServerPool manages the list of backends and the current status.
type ServerPool struct {
	backends []*url.URL
	current  uint64
}

// AddBackend adds a new backend URL to the pool.
func (s *ServerPool) AddBackend(backendUrl string) {
	u, err := url.Parse(backendUrl)
	if err != nil {
		log.Fatal(err)
	}
	s.backends = append(s.backends, u)
}

// GetNextBackend returns the next backend to serve a request using Round Robin.
func (s *ServerPool) GetNextBackend() *url.URL {
	// Atomically increment the counter to ensure thread safety.
	// We use modulo to wrap around the list of backends.
	next := atomic.AddUint64(&s.current, 1)
	index := int(next % uint64(len(s.backends)))
	return s.backends[index]
}

// lbHandler is the main HTTP handler for our Load Balancer.
func lbHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Choose a backend
	target := serverPool.GetNextBackend()

	// 2. Create a Reverse Proxy
	// This standard library tool will forward the request to the target
	// and send the response back to the client.
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Optional: Log the routing decision
	fmt.Printf("LB: Balancing request to %s\n", target)

	// 3. Serve the request using the proxy
	// The Update of the request URL to the target happens inside ServeHTTP
	proxy.ServeHTTP(w, r)
}

var serverPool ServerPool

func main() {
	// Define our backend servers (these must match the ports we run our scaling app on)
	serverPool.AddBackend("http://localhost:8081")
	serverPool.AddBackend("http://localhost:8082")
	serverPool.AddBackend("http://localhost:8083")

	// Start the Load Balancer on port 8000
	port := ":8000"
	fmt.Printf("Load Balancer started on port %s\n", port)
	fmt.Println("Forwarding traffic to:")
	for _, b := range serverPool.backends {
		fmt.Printf(" - %s\n", b)
	}

	// Route all traffic to our load balancer handler
	http.HandleFunc("/", lbHandler)
	log.Fatal(http.ListenAndServe(port, nil))
}

