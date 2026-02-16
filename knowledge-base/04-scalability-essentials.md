# Scalability Essentials: From Load Balancers to Sharding

**Target Audience:** Engineers who understand code but want to understand *architecture*. This guide explains how to take a single server and scale it to handle millions of users.

---

## 1. Vertical vs. Horizontal Scaling

### The Concept
*   **Vertical Scaling (Scale Up):** Buying a bigger computer.
    *   **Pros:** Simple. No code changes.
    *   **Cons:** Expensive. Has a hard limit (you can't buy a CPU with 10,000 cores). Single point of failure.
*   **Horizontal Scaling (Scale Out):** Buying *more* computers.
    *   **Pros:** Infinite scale. Cheaper commodity hardware. High availability.
    *   **Cons:** Complex. Requires Load Balancing, Distributed Data, and Network calls.

### Simulation: The "Stateless" Backend
To scale horizontally, your application must be **Stateless**. It cannot store user sessions in local memory, because the next request might hit a different server.

```go
package main

import (
    "flag"
    "fmt"
    "net/http"
)

func main() {
    // We can run multiple instances of this app on different ports.
    // go run main.go -port=8081
    // go run main.go -port=8082
    port := flag.String("port", "8080", "Port to run on")
    flag.Parse()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // The server identifies itself, proving which instance handled the request.
        fmt.Fprintf(w, "Hello from Server running on port %s\n", *port)
    })

    http.ListenAndServe(":"+*port, nil)
}
```

---

## 2. Load Balancing (L4 vs L7)

### The Concept
A Load Balancer (LB) sits in front of your servers and distributes traffic.

*   **L4 (Layer 4 - Transport):** Routes based on IP and Port.
    *   **Pros:** Extremely fast.
    *   **Cons:** "Dumb". Cannot see the URL or Headers. Can't route `/api` to Server A and `/images` to Server B.
*   **L7 (Layer 7 - Application):** Routes based on HTTP content (URL, Cookies, Headers).
    *   **Pros:** Smart routing (Microservices). Can terminate SSL.
    *   **Cons:** Slower (needs to decrypt/inspect packets).

### Implementation: A Simple Round-Robin L7 LB
This Go code acts as a reverse proxy, forwarding requests to backend servers in a loop.

```go
package main

import (
    "net/http"
    "net/http/httputil"
    "net/url"
    "sync/atomic"
)

type LoadBalancer struct {
    backends []*url.URL
    counter  uint64
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Round Robin Algorithm
    // Atomic increment ensures thread safety without heavy locks.
    current := atomic.AddUint64(&lb.counter, 1)
    index := current % uint64(len(lb.backends))
    
    target := lb.backends[index]
    
    // Reverse Proxy: Forwards request to target, sends response back to client.
    proxy := httputil.NewSingleHostReverseProxy(target)
    proxy.ServeHTTP(w, r)
}

func main() {
    u1, _ := url.Parse("http://localhost:8081")
    u2, _ := url.Parse("http://localhost:8082")

    lb := &LoadBalancer{backends: []*url.URL{u1, u2}}
    http.ListenAndServe(":8000", lb)
}
```

---

## 3. Caching Strategies

### The Concept
Reading from memory (RAM) is nanoseconds. Reading from Disk/Network is milliseconds. Caching saves time.

*   **Cache-Aside (Lazy Loading):**
    1.  App asks Cache.
    2.  If Miss, App asks DB.
    3.  App writes DB result to Cache.
    *   **Pros:** Only caches what is needed. Resilient to cache failure.
    *   **Cons:** First request is slow (thundering herd risk).
*   **Write-Through:**
    1.  App writes to Cache AND DB simultaneously.
    *   **Pros:** Data is always consistent.
    *   **Cons:** Slow writes. Caches data that might never be read.

### Implementation: Cache-Aside Pattern

```go
package main

import "fmt"

type Database map[string]string
type Cache map[string]string

type App struct {
    DB    Database
    Cache Cache
}

func (a *App) GetUser(id string) string {
    // 1. Check Cache
    if val, ok := a.Cache[id]; ok {
        fmt.Println("Cache Hit")
        return val
    }

    // 2. Fetch from DB
    fmt.Println("Cache Miss -> Fetching from DB")
    val := a.DB[id]

    // 3. Update Cache
    a.Cache[id] = val
    return val
}

func main() {
    app := App{
        DB:    map[string]string{"1": "Alice"},
        Cache: make(map[string]string),
    }

    app.GetUser("1") // Miss
    app.GetUser("1") // Hit
}
```

---

## 4. Database Sharding

### The Concept
When a database gets too big for one disk (e.g., 10TB), you split it across multiple servers.

*   **Vertical Partitioning:** Columns. (e.g., "User Table" on Server A, "Orders Table" on Server B).
*   **Horizontal Partitioning (Sharding):** Rows. (e.g., Users A-M on Server A, Users N-Z on Server B).

**The Challenge:** How do you know which server has "Alice"?
**The Solution:** Consistent Hashing (covered in Advanced Module) or simple Modulo Hashing.

### Implementation: Modulo Sharding
`ShardID = Hash(Key) % NumShards`

```go
package main

import (
    "fmt"
    "hash/crc32"
)

type Shard struct {
    ID    int
    Store map[string]string
}

type ShardedDB struct {
    Shards []*Shard
}

func (db *ShardedDB) GetShard(key string) *Shard {
    // CRC32 is a fast hashing algorithm.
    hash := crc32.ChecksumIEEE([]byte(key))
    // Modulo arithmetic determines the bucket.
    index := int(hash) % len(db.Shards)
    return db.Shards[index]
}

func (db *ShardedDB) Put(key, value string) {
    shard := db.GetShard(key)
    fmt.Printf("Saving '%s' to Shard %d\n", key, shard.ID)
    shard.Store[key] = value
}

func main() {
    // Create 3 Shards
    db := &ShardedDB{
        Shards: []*Shard{
            {ID: 0, Store: make(map[string]string)},
            {ID: 1, Store: make(map[string]string)},
            {ID: 2, Store: make(map[string]string)},
        },
    }

    // Deterministic distribution
    db.Put("Alice", "Data")   // Shard 1
    db.Put("Bob", "Data")     // Shard 2
    db.Put("Charlie", "Data") // Shard 1
}
```

---

## 5. Content Delivery Networks (CDN)

### The Concept
Speed of light is finite. If your server is in New York and the user is in Sydney, latency is ~200ms minimum.
A **CDN** is a network of servers ("Edge Nodes") distributed globally.

1.  User requests `logo.png`.
2.  Request hits nearest Edge Node (Sydney).
3.  If Edge has it -> Return instantly.
4.  If Edge misses -> Fetch from Origin (NY), cache it in Sydney, return it.

### Implementation: CDN Simulation

```go
package main

import "fmt"

type Origin struct{}
func (o *Origin) Fetch(file string) string {
    fmt.Println("   [Origin] Reading from disk (Slow)...")
    return "FileContent"
}

type EdgeNode struct {
    Origin *Origin
    Cache  map[string]string
}

func (e *EdgeNode) Get(file string) string {
    if val, ok := e.Cache[file]; ok {
        fmt.Println("[Edge] Cache Hit! (Fast)")
        return val
    }
    
    fmt.Println("[Edge] Cache Miss. Contacting Origin...")
    val := e.Origin.Fetch(file)
    e.Cache[file] = val
    return val
}

func main() {
    cdn := &EdgeNode{
        Origin: &Origin{},
        Cache:  make(map[string]string),
    }

    cdn.Get("video.mp4") // Slow
    cdn.Get("video.mp4") // Fast
}
```
