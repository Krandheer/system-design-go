package main

import (
	"fmt"
	"time"
)

// Server is our complex object. Notice it has many fields.
// Creating this with a single constructor would be messy.
type Server struct {
	Host    string
	Port    int
	Timeout time.Duration
	UseTLS  bool
	MaxConn int
}

// --- The Builder ---

// ServerBuilder is the builder type for our Server.
type ServerBuilder struct {
	server *Server
}

// NewServerBuilder creates and returns a new builder instance.
// It initializes the Server with some sensible defaults.
func NewServerBuilder(host string) *ServerBuilder {
	return &ServerBuilder{
		server: &Server{
			Host:    host,
			Port:    8080,             // Default port
			Timeout: 30 * time.Second, // Default timeout
			UseTLS:  false,            // Default TLS
			MaxConn: 100,              // Default max connections
		},
	}
}

// --- Fluent Methods ---
// Each of these methods modifies the server configuration and returns
// the builder pointer, allowing for method chaining.

func (b *ServerBuilder) WithPort(port int) *ServerBuilder {
	b.server.Port = port
	return b
}

func (b *ServerBuilder) WithTimeout(timeout time.Duration) *ServerBuilder {
	b.server.Timeout = timeout
	return b
}

func (b *ServerBuilder) WithTLS(useTLS bool) *ServerBuilder {
	b.server.UseTLS = useTLS
	return b
}

func (b *ServerBuilder) WithMaxConn(maxConn int) *ServerBuilder {
	b.server.MaxConn = maxConn
	return b
}

// --- The Final Build Step ---

// Build returns the fully constructed Server object.
// This is the final step in the builder process.
func (b *ServerBuilder) Build() *Server {
	return b.server
}

func main() {
	fmt.Println("--- Building a simple HTTP server ---")
	// We start with a required parameter (host) and then chain optional settings.
	// This is much more readable than a constructor with many arguments.
	httpServer := NewServerBuilder("localhost").
		WithPort(80).
		Build()

	fmt.Printf("HTTP Server Config: %+v\n\n", *httpServer)

	fmt.Println("--- Building a complex HTTPS server ---")
	// Here's a more complex example with more options chained.
	httpsServer := NewServerBuilder("api.example.com").
		WithPort(443).
		WithTLS(true).
		WithTimeout(10 * time.Second).
		WithMaxConn(500).
		Build()

	fmt.Printf("HTTPS Server Config: %+v\n", *httpsServer)
}
