# Reliability & APIs: CAP, Consistency, and Interface Design

**Target Audience:** Engineers who need to make hard choices about data integrity and system availability, and who need to design clean APIs for their services.

---

## 1. The CAP Theorem

### The Concept
In a distributed system (where data is on multiple nodes), you can only have **two** of the following three:
*   **Consistency (C):** Every read receives the most recent write or an error. (All nodes see the same data at the same time).
*   **Availability (A):** Every request receives a (non-error) response, without the guarantee that it contains the most recent write.
*   **Partition Tolerance (P):** The system continues to operate despite an arbitrary number of messages being dropped or delayed by the network between nodes.

**The Reality:** In a distributed system, **P is mandatory**. Networks fail. Cables get cut.
Therefore, you must choose between **CP** (Consistency) and **AP** (Availability).

### Simulation: The Hard Choice

```go
package main

import "fmt"

type SystemMode int
const (
    CP SystemMode = iota // Consistency First (e.g., Banking)
    AP                   // Availability First (e.g., Twitter Feed)
)

type Node struct {
    Data string
}

type Cluster struct {
    Mode      SystemMode
    Master    *Node
    Replica   *Node
    Partition bool // Is the network broken?
}

func (c *Cluster) Write(data string) {
    fmt.Printf("Writing '%s'...\n", data)
    
    if c.Partition {
        if c.Mode == CP {
            fmt.Println("Error: Network Partition. Cannot replicate. Write REJECTED.")
            return
        }
        if c.Mode == AP {
            fmt.Println("Warning: Network Partition. Writing to Master only.")
            c.Master.Data = data
            return
        }
    }

    c.Master.Data = data
    c.Replica.Data = data
    fmt.Println("Success: Written to both nodes.")
}

func main() {
    // 1. CP System (Bank)
    fmt.Println("--- CP System (Bank) ---")
    bank := &Cluster{Mode: CP, Master: &Node{}, Replica: &Node{}, Partition: true}
    bank.Write("Balance: $100") // Fails. Better to fail than have different balances.

    // 2. AP System (Social)
    fmt.Println("\n--- AP System (Social) ---")
    social := &Cluster{Mode: AP, Master: &Node{}, Replica: &Node{}, Partition: true}
    social.Write("New Tweet") // Succeeds. Replica is stale, but Master is updated.
}
```

---

## 2. Consistency Patterns

### The Concept
If you choose AP (Availability), your data will be inconsistent for a while. How do you handle that?

*   **Strong Consistency:** You wait for ALL replicas to acknowledge the write before telling the user "Success". (Slow, CP).
*   **Eventual Consistency:** You write to one node and return "Success". The data propagates to others in the background. (Fast, AP).

### Implementation: Eventual Consistency (Replication Lag)

```go
package main

import (
    "fmt"
    "time"
)

type Database struct {
    Value string
}

func main() {
    master := &Database{Value: "v1"}
    replica := &Database{Value: "v1"}

    // User writes to Master
    fmt.Println("Write: v2 to Master")
    master.Value = "v2"

    // Replication happens asynchronously
    go func() {
        time.Sleep(2 * time.Second) // Lag
        replica.Value = master.Value
        fmt.Println("\n[Background] Replication Complete.")
    }()

    // User reads immediately from Replica
    fmt.Printf("Read from Replica: %s (Stale!)\n", replica.Value)

    time.Sleep(3 * time.Second)
    fmt.Printf("Read from Replica: %s (Consistent)\n", replica.Value)
}
```

---

## 3. Availability Patterns (Failover)

### The Concept
How do you ensure your system stays up when a server crashes?

*   **Active-Passive:** Server A handles all traffic. Server B sits idle. If A dies, B takes over. (Simple, but wastes resources).
*   **Active-Active:** Both A and B handle traffic. If A dies, B handles 100% of the load. (Efficient, but B must be powerful enough).

### Implementation: Active-Passive Failover

```go
package main

import "fmt"

type Server struct {
    ID    string
    Alive bool
}

type LoadBalancer struct {
    Active  *Server
    Passive *Server
}

func (lb *LoadBalancer) GetServer() *Server {
    if lb.Active.Alive {
        return lb.Active
    }
    fmt.Println("ALERT: Active is DOWN. Switching to Passive.")
    return lb.Passive
}

func main() {
    s1 := &Server{ID: "Primary", Alive: true}
    s2 := &Server{ID: "Backup", Alive: true}
    lb := &LoadBalancer{Active: s1, Passive: s2}

    fmt.Println("Request 1 served by:", lb.GetServer().ID)

    // Crash!
    s1.Alive = false
    
    fmt.Println("Request 2 served by:", lb.GetServer().ID)
}
```

---

## 4. API Design: REST vs. GraphQL vs. gRPC

### The Comparison

| Feature | REST | GraphQL | gRPC |
| :--- | :--- | :--- | :--- |
| **Protocol** | HTTP/1.1 (Text) | HTTP/1.1 (Text) | HTTP/2 (Binary) |
| **Data Format** | JSON | JSON | Protobuf |
| **Philosophy** | Resources (`/users/1`) | Query Language | Function Calls |
| **Pros** | Simple, Cacheable, Universal. | No Over-fetching, Flexible. | Extremely Fast, Typed. |
| **Cons** | Over-fetching (getting too much data). | Complex Caching, N+1 problem. | Browser support requires proxy. |

### Implementation: The Three Flavors

```go
package main

import (
    "encoding/json"
    "fmt"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"` // Heavy field we might not want
}

var dbUser = User{ID: 1, Name: "Alice", Email: "a@a.com", Age: 30}

// 1. REST: Returns everything.
func HandleREST() string {
    b, _ := json.Marshal(dbUser)
    return string(b)
}

// 2. GraphQL: Returns only what is asked.
func HandleGraphQL(fields []string) string {
    res := make(map[string]interface{})
    for _, f := range fields {
        if f == "name" { res["name"] = dbUser.Name }
        if f == "email" { res["email"] = dbUser.Email }
    }
    b, _ := json.Marshal(res)
    return string(b)
}

// 3. gRPC: Strict Types (Simulated)
type UserResponse struct {
    Name string
}
func HandleGRPC() UserResponse {
    return UserResponse{Name: dbUser.Name}
}

func main() {
    fmt.Println("REST:", HandleREST()) 
    // Output: {"id":1,"name":"Alice","email":"a@a.com","age":30}

    fmt.Println("GraphQL:", HandleGraphQL([]string{"name"})) 
    // Output: {"name":"Alice"}

    fmt.Printf("gRPC: %+v\n", HandleGRPC()) 
    // Output: {Name:Alice}
}
```
