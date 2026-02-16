package main

import (
	"fmt"
	"sync"
	"time"
)

// Server simulates a backend server.
type Server struct {
	ID      string
	IsAlive bool
	mu      sync.RWMutex
}

func (s *Server) Ping() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.IsAlive
}

func (s *Server) Kill() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IsAlive = false
	fmt.Printf("--- ALERT: %s has CRASHED! ---\n", s.ID)
}

// LoadBalancer manages the traffic and failover.
type LoadBalancer struct {
	Active  *Server
	Passive *Server
	mu      sync.Mutex
}

// Serve handles a client request.
func (lb *LoadBalancer) Serve() {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	fmt.Printf("Request served by: %s\n", lb.Active.ID)
}

// HealthCheck monitors the active server and performs failover.
func (lb *LoadBalancer) StartHealthCheck() {
	for {
		time.Sleep(500 * time.Millisecond)
		
		lb.mu.Lock()
		currentServer := lb.Active
		lb.mu.Unlock()

		if !currentServer.Ping() {
			fmt.Println("Health Check: Active server is DOWN. Initiating Failover...")
			lb.Failover()
		} else {
			// fmt.Println("Health Check: Active server is OK.")
		}
	}
}

func (lb *LoadBalancer) Failover() {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// Swap Active and Passive
	// In a real scenario, you'd check if Passive is alive first.
	fmt.Printf("FAILOVER: Promoting %s to Active.\n", lb.Passive.ID)
	lb.Active = lb.Passive
	// The old active is now the passive (and dead) server.
}

func main() {
	serverA := &Server{ID: "Server A (Primary)", IsAlive: true}
	serverB := &Server{ID: "Server B (Backup)", IsAlive: true}

	lb := &LoadBalancer{
		Active:  serverA,
		Passive: serverB,
	}

	// Start the background health monitor
	go lb.StartHealthCheck()

	// Simulate traffic
	for i := 1; i <= 10; i++ {
		lb.Serve()
		time.Sleep(300 * time.Millisecond)

		// Simulate a crash at request #4
		if i == 4 {
			serverA.Kill()
		}
	}
}
