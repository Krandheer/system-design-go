package main

import (
	"fmt"
	"hash/crc32"
)

// Record represents a piece of data in our DB.
type Record struct {
	Key   string
	Value string
}

// Shard represents a single database instance (e.g., a PostgreSQL server).
type Shard struct {
	ID    int
	Store map[string]Record
}

func NewShard(id int) *Shard {
	return &Shard{
		ID:    id,
		Store: make(map[string]Record),
	}
}

// ShardedDatabase manages a collection of Shards.
type ShardedDatabase struct {
	Shards []*Shard
}

func NewShardedDatabase(numShards int) *ShardedDatabase {
	shards := make([]*Shard, numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = NewShard(i)
	}
	return &ShardedDatabase{Shards: shards}
}

// getShardIndex determines which shard handles a given key.
// It uses a Hash function (CRC32) and Modulo arithmetic.
// Hash(Key) % NumShards = ShardIndex
func (sdb *ShardedDatabase) getShardIndex(key string) int {
	hash := crc32.ChecksumIEEE([]byte(key))
	return int(hash) % len(sdb.Shards)
}

// Save stores data in the correct shard.
func (sdb *ShardedDatabase) Save(key, value string) {
	shardIndex := sdb.getShardIndex(key)
	shard := sdb.Shards[shardIndex]

	fmt.Printf("Saving key '%s' -> Shard %d\n", key, shard.ID)
	shard.Store[key] = Record{Key: key, Value: value}
}

// Get retrieves data from the correct shard.
func (sdb *ShardedDatabase) Get(key string) (string, bool) {
	shardIndex := sdb.getShardIndex(key)
	shard := sdb.Shards[shardIndex]

	record, ok := shard.Store[key]
	return record.Value, ok
}

func main() {
	// 1. Initialize a DB with 3 Shards
	db := NewShardedDatabase(3)

	fmt.Println("--- Storing Data ---")
	// Notice how different users get mapped to different shards based on their Key
	users := []string{"Alice", "Bob", "Charlie", "Dave", "Eve", "Frank", "Grace", "Heidi", "Ivan", "Judy", "Mallory"}
	for _, user := range users {
		db.Save(user, "Data for "+user)
	}

	fmt.Println("\n--- Retrieving Data ---")
	// When we read back, we use the same math to find where the data lives.
	val, _ := db.Get("Alice")
	fmt.Printf("Got Alice: %s\n", val)
	
	val, _ = db.Get("Bob")
	fmt.Printf("Got Bob: %s\n", val)
}

