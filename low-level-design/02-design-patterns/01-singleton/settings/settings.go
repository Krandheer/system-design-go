// Package settings provides a singleton for application configuration.
package settings

import (
	"fmt"
	"sync"
)

// Settings defines the structure of our configuration.
// It's unexported to prevent direct instantiation from other packages.
type settings struct {
	DatabaseURL string
	APIKey      string
	Port        int
}

var (
	// instance is the pointer to the single instance of the settings.
	// It's unexported.
	instance *settings

	// once is a sync.Once object that will ensure the initialization
	// code is executed exactly once.
	once sync.Once
)

// GetInstance is the public, global access point to the singleton.
// This is the only function that should be called from outside this package.
func GetInstance() *settings {
	// The Do method of sync.Once takes a function as an argument.
	// This function will be executed only on the very first call to Do.
	// All subsequent calls to Do will do nothing, but will block until
	// the first call's function has returned. This makes it thread-safe.
	once.Do(func() {
		fmt.Println("Initializing settings for the first and only time...")
		// Here, we create the single instance.
		// In a real application, you might load these values from a file
		// or environment variables.
		instance = &settings{
			DatabaseURL: "postgres://user:pass@localhost/db",
			APIKey:      "super-secret-key",
			Port:        8080,
		}
	})
	return instance
}

// GetPort is an example of a method on our singleton struct.
func (s *settings) GetPort() int {
	return s.Port
}
