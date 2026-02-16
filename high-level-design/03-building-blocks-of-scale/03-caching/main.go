package main

import (
	"fmt"
	"sync"
	"time"
)

// User represents the data we want to retrieve.
type User struct {
	ID    string
	Name  string
	Email string
}

// Database simulates a slow persistent storage.
type Database struct {
	data map[string]User
}

func NewDatabase() *Database {
	return &Database{
		data: map[string]User{
			"1": {ID: "1", Name: "Alice", Email: "alice@example.com"},
			"2": {ID: "2", Name: "Bob", Email: "bob@example.com"},
		},
	}
}

// GetUser simulates a slow database query.
func (db *Database) GetUser(id string) (User, bool) {
	time.Sleep(2 * time.Second) // Simulate network/disk latency
	user, ok := db.data[id]
	return user, ok
}

// Cache simulates a fast in-memory store (like Redis).
type Cache struct {
	store map[string]User
	mutex sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]User),
	}
}

func (c *Cache) Get(id string) (User, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	user, ok := c.store[id]
	return user, ok
}

func (c *Cache) Set(id string, user User) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.store[id] = user
}

// Application represents our backend service.
type Application struct {
	db    *Database
	cache *Cache
}

// GetProfile implements the Cache-Aside pattern.
func (app *Application) GetProfile(userID string) User {
	fmt.Printf("Requesting profile for user %s...\n", userID)

	// 1. Check Cache
	user, found := app.cache.Get(userID)
	if found {
		fmt.Println(" -> Cache HIT! Returning fast.")
		return user
	}

	fmt.Println(" -> Cache MISS. Fetching from DB (slow)...")

	// 2. Fetch from DB
	user, found = app.db.GetUser(userID)
	if !found {
		fmt.Println(" -> User not found in DB.")
		return User{}
	}

	// 3. Update Cache
	fmt.Println(" -> Writing to Cache...")
	app.cache.Set(userID, user)

	return user
}

func main() {
	app := &Application{
		db:    NewDatabase(),
		cache: NewCache(),
	}

	// First Request: Cache Miss
	start := time.Now()
	app.GetProfile("1")
	fmt.Printf("First call took: %v\n\n", time.Since(start))

	// Second Request: Cache Hit
	start = time.Now()
	app.GetProfile("1")
	fmt.Printf("Second call took: %v\n", time.Since(start))
}

