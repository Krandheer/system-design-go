# Advanced Distributed Algorithms: The "Expert" Level

**Target Audience:** Staff+ Engineers and Architects. This guide covers the complex algorithms that power the infrastructure of the internet (DynamoDB, Cassandra, Kafka, Kubernetes).

---

## 1. Consistent Hashing (The Ring)

### The Problem
Traditional sharding uses `Hash(key) % N`.
If you add/remove a server (`N` changes), **ALL** keys are remapped. This causes a cache stampede or database downtime.

### The Solution
Map both Servers and Keys to a circle (0-360°).
*   **Placement:** Hash the Server IP to a point on the circle.
*   **Assignment:** Hash the Key to a point. Walk **clockwise** until you find a server.
*   **Stability:** If a server is removed, only the keys belonging to *that* server are moved to the next neighbor. All other keys stay put.

### Implementation: The Ring with Virtual Nodes

```go
package main

import (
    "fmt"
    "hash/crc32"
    "sort"
    "strconv"
)

type HashRing struct {
    keys     []int          // Sorted hash values
    hashMap  map[int]string // Hash -> Node Name
    replicas int            // Virtual nodes per physical node
}

func NewRing(replicas int) *HashRing {
    return &HashRing{replicas: replicas, hashMap: make(map[int]string)}
}

func (h *HashRing) Add(node string) {
    for i := 0; i < h.replicas; i++ {
        // Virtual Node: "NodeA#1", "NodeA#2"
        hash := int(crc32.ChecksumIEEE([]byte(node + strconv.Itoa(i))))
        h.keys = append(h.keys, hash)
        h.hashMap[hash] = node
    }
    sort.Ints(h.keys)
}

func (h *HashRing) Get(key string) string {
    hash := int(crc32.ChecksumIEEE([]byte(key)))
    
    // Binary Search: Find the first server clockwise
    idx := sort.Search(len(h.keys), func(i int) bool {
        return h.keys[i] >= hash
    })

    if idx == len(h.keys) {
        idx = 0 // Wrap around
    }

    return h.hashMap[h.keys[idx]]
}

func main() {
    ring := NewRing(3)
    ring.Add("Server-A")
    ring.Add("Server-B")

    fmt.Println("User1 ->", ring.Get("User1"))
    fmt.Println("User2 ->", ring.Get("User2"))
}
```

---

## 2. Distributed ID Generation (Snowflake)

### The Problem
In distributed systems, you cannot use `AUTO_INCREMENT` (requires a central DB lock).
UUIDs (128-bit) are unique but unordered (bad for DB indexing) and large.

### The Solution: Twitter Snowflake
A 64-bit integer composed of:
`[ 1 bit Sign | 41 bits Timestamp | 10 bits Machine ID | 12 bits Sequence ]`

*   **Sortable:** Time is the most significant part.
*   **Unique:** Machine ID prevents conflicts between servers. Sequence prevents conflicts on the same server (up to 4096/ms).

### Implementation: Bitwise Generation

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

const (
    epoch = 1609459200000 // Custom Epoch (Jan 1 2021)
    machineBits  = 10
    sequenceBits = 12
)

type Snowflake struct {
    mu        sync.Mutex
    lastTime  int64
    machineID int64
    sequence  int64
}

func (s *Snowflake) NextID() int64 {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now().UnixMilli()
    if now == s.lastTime {
        s.sequence++
        // Overflow handling omitted for brevity
    } else {
        s.sequence = 0
    }
    s.lastTime = now

    // Shift and OR (|) to pack bits
    return ((now - epoch) << (machineBits + sequenceBits)) |
           (s.machineID << sequenceBits) |
           s.sequence
}

func main() {
    node := &Snowflake{machineID: 1}
    fmt.Println(node.NextID())
    fmt.Println(node.NextID()) // Slightly larger
}
```

---

## 3. Bloom Filters

### The Problem
Checking if a username exists in a DB of 1 Billion users is slow (Disk I/O).
Checking a Cache is faster, but RAM is expensive (1 Billion strings = 100GB+).

### The Solution
A probabilistic data structure that uses **Bits**, not Strings.
*   "Definitely NOT in set" (100% accurate).
*   "PROBABLY in set" (Small false positive rate).

### Implementation: Double Hashing

```go
package main

import (
    "fmt"
    "hash/fnv"
)

type BloomFilter struct {
    bitSet []bool
    size   uint32
}

func (bf *BloomFilter) Add(data string) {
    h1, h2 := bf.hash(data)
    bf.bitSet[h1%bf.size] = true
    bf.bitSet[h2%bf.size] = true
}

func (bf *BloomFilter) Check(data string) bool {
    h1, h2 := bf.hash(data)
    // If ANY bit is 0, it's definitely missing.
    if !bf.bitSet[h1%bf.size] || !bf.bitSet[h2%bf.size] {
        return false
    }
    return true // Probably present
}

func (bf *BloomFilter) hash(s string) (uint32, uint32) {
    h := fnv.New32a()
    h.Write([]byte(s))
    v1 := h.Sum32()
    return v1, v1 + 17 // Simulate second hash
}

func main() {
    bf := &BloomFilter{bitSet: make([]bool, 100), size: 100}
    bf.Add("apple")
    
    fmt.Println("apple?", bf.Check("apple")) // True
    fmt.Println("car?", bf.Check("car"))     // False
}
```

---

## 4. Distributed Consensus (Raft)

### The Problem
How do you get 3 servers to agree on "Who is the leader?" or "What is the value of X?" when networks fail and servers crash?

### The Solution: Quorum
*   **Quorum:** `(N / 2) + 1`. (e.g., 3 nodes -> need 2 votes).
*   **Raft:**
    1.  **Leader Election:** Nodes start as Followers. If they don't hear from a Leader, they become Candidates and ask for votes. If they get a Quorum, they become Leader.
    2.  **Log Replication:** The Leader accepts writes and replicates them. It only commits if a Quorum acknowledges the write.

### Implementation: Leader Election State Machine

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

type State int
const (
    Follower State = iota
    Candidate
    Leader
)

type Node struct {
    ID    int
    State State
    Votes int
}

func (n *Node) ElectionTimeout() {
    // If timeout fires, start election
    n.State = Candidate
    n.Votes = 1 // Vote for self
    fmt.Printf("[%d] Timeout! Becoming Candidate.\n", n.ID)
    n.RequestVotes()
}

func (n *Node) RequestVotes() {
    // Simulate getting a vote from a peer
    // In real Raft, this is an RPC call
    n.Votes++ 
    if n.Votes >= 2 { // Quorum for cluster of 3 is 2
        n.State = Leader
        fmt.Printf("[%d] Quorum reached! Becoming LEADER.\n", n.ID)
        n.SendHeartbeats()
    }
}

func (n *Node) SendHeartbeats() {
    fmt.Printf("[%d] Sending Heartbeats (I am Leader)...\n", n.ID)
}

func main() {
    node := &Node{ID: 1, State: Follower}
    
    // Simulate time passing without leader
    time.Sleep(1 * time.Second)
    node.ElectionTimeout()
}
```
