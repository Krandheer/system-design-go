# System Design and Golang Mastery: A Learning Plan

This document is our "living" plan to track progress, store notes, and maintain context throughout our system design journey.

---

## Part 1: Low-Level Design (LLD) with Go - The Idiomatic Approach

This part directly addresses how to achieve OOP principles in a non-OOP language, focusing on Go's unique strengths.

*   **Module 1: Thinking in Go**
    *   [x] Structs and Methods vs. Classes
    *   [x] Interfaces: The Power of Implicit Contracts (Polymorphism)
    *   [x] Composition over Inheritance: Go's core philosophy
    *   [x] Encapsulation through Packages
    *   [x] Error Handling Strategies

*   **Module 2: Foundational Design Patterns in Go**
    *   **Creational Patterns**
        *   [x] Singleton
        *   [x] Factory
        *   [x] Builder
    *   **Structural Patterns**
        *   [x] Adapter
        *   [x] Decorator
        *   [x] Facade
    *   **Behavioral Patterns**
        *   [x] Strategy
        *   [x] Observer
        *   [x] Chain of Responsibility
    *   **Concurrency Patterns**
        *   [x] Worker Pools
        *   [x] Fan-in / Fan-out
        *   [x] Rate Limiting (LLD Implementation)

---

## Part 2: High-Level Design (HLD) - The Big Picture

*   **Module 3: The Building Blocks of Scale**
    *   [x] Vertical vs. Horizontal Scaling
    *   [x] Load Balancers (L4 vs. L7)
    *   [x] Caching: Strategies and Patterns (e.g., Write-Through, Write-Around, Read-Aside)
    *   [x] Databases: SQL vs. NoSQL, Sharding, Replication, Indexes
    *   [x] Asynchronism: Message Queues and Event-Driven Architecture
    *   [x] Content Delivery Networks (CDN)
    *   [x] Proxies (Forward vs. Reverse)

*   **Module 4: Ensuring Reliability and Consistency**
    *   [x] CAP Theorem explained practically
    *   [x] Consistency Patterns
    *   [x] Availability Patterns (Failover, Redundancy)
    *   [x] API Design: REST vs. GraphQL vs. gRPC

*   **Module 5: Advanced Distributed Concepts**
    *   [x] Consistent Hashing (Ring & Virtual Nodes)
    *   [x] Distributed ID Generation (Snowflake, UUID)
    *   [x] Distributed Transactions (Saga Pattern, 2PC)
    *   [x] Bloom Filters (Probabilistic Data Structures)
    *   [x] Quorum & Consensus (Raft, Paxos basics)

---

## Part 3: System Design Case Studies (HLD + LLD Implementation)

We will tackle these classic interview problems from requirements gathering to HLD, LLD, and finally, implementing core components in Go.

*   [ ] **Case Study 1:** URL Shortener (like TinyURL)
*   [ ] **Case Study 2:** Web Crawler
*   [ ] **Case Study 3:** Distributed Task Queue / Job Scheduler
*   [ ] **Case Study 4:** Pastebin / Codebin service
*   [ ] **Case Study 5:** Rate Limiter (Distributed)
*   [ ] **Case Study 6:** Search Autocomplete System (Trie-based)
*   [ ] **Case Study 7:** Ride-Sharing App (like Uber)
*   [ ] **Case Study 8:** Chat Application (like WhatsApp)
*   [ ] **Case Study 9:** Social Media Feed and News Feed (like Twitter/Facebook)
*   [ ] **Case Study 10:** Video Streaming Service (like YouTube/Netflix)
*   [ ] **Case Study 11:** Distributed Key-Value Store (like DynamoDB)

---

## Session Log & Notes

*(We will add summaries of our sessions here to maintain context)*

---
**Session 1: Module 1 - Thinking in Go**

*   **Summary:** We established the foundational principles of idiomatic Go, focusing on how it achieves object-oriented concepts without traditional classes.
*   **Key Concepts Covered:**
    *   **Structs & Methods:** Go separates data (structs) from behavior (methods), unlike classes which bundle them.
    *   **Interfaces:** Learned about implicit satisfaction ("duck typing"), enabling polymorphism and decoupled code. Implemented a `Shape` interface with `Rectangle` and `Circle`.
    *   **Composition over Inheritance:** Used struct embedding to build complex types from simpler ones (e.g., a `Manager` "has-a" `Worker`), promoting flexibility over rigid hierarchies.
    *   **Encapsulation via Packages:** Understood that visibility is controlled by capitalization (e.g., `Public` vs. `private`). Created an `employee` package with a private `salary` field accessed via public methods.
    *   **Error Handling:** Covered Go's core philosophy that "errors are values." We practiced explicit `if err != nil` checks, creating custom error types, and wrapping errors to add context.
*   **Code Location:** `low-level-design/01-thinking-in-go/`
---
**Session 2: Module 2 - Creational Design Patterns**

*   **Summary:** We implemented the three foundational creational patterns in idiomatic Go. The focus was on managing object creation in a clean, flexible, and concurrency-safe way.
*   **Key Concepts Covered:**
    *   **Singleton:** Ensured a struct has only one instance and provided a global access point. We used `sync.Once` to make our singleton thread-safe.
    *   **Factory:** Decoupled object creation from client code. We created a `GetNotifier` factory that returned a `Notifier` interface, hiding the concrete `EmailNotifier` and `SMSNotifier` types.
    *   **Builder:** Constructed complex objects step-by-step. We built a `ServerBuilder` with a fluent interface to create `Server` objects with multiple optional configurations in a readable way.
*   **Code Location:** `low-level-design/02-design-patterns/`
---
**Session 3: Module 2 - Structural Design Patterns**

*   **Summary:** We explored patterns that deal with the composition of types and objects to form larger, flexible structures. The focus was on making different parts of a system work together seamlessly.
*   **Key Concepts Covered:**
    *   **Adapter:** Made incompatible interfaces compatible. We created a `StripeAdapter` to allow our system's `PaymentProcessor` interface to work with an external `StripeClient`.
    *   **Decorator:** Attached new behaviors to objects dynamically. We wrapped a base `PlainPizza` with `CheeseTopping` and `TomatoTopping` decorators to build up its price and description at runtime.
    *   **Facade:** Provided a simplified interface to a complex subsystem. We built a `HomeTheaterFacade` with simple `WatchMovie()` and `EndMovie()` methods to hide the complexity of managing individual components like projectors and amplifiers.
*   **Code Location:** `low-level-design/02-design-patterns/`
---
**Session 4: Module 2 - Behavioral & Concurrency Design Patterns**

*   **Summary:** We concluded our study of foundational patterns by focusing on how objects communicate and how to manage concurrent operations. We covered patterns for encapsulating algorithms, notifying objects of state changes, processing requests along a chain, and managing concurrent execution.
*   **Key Concepts Covered:**
    *   **Strategy:** Made a family of algorithms interchangeable (`Navigator` with different `RouteStrategy` implementations).
    *   **Observer:** Created a subscription mechanism for one-to-many notifications (`NewsAgency` notifying `NewsSubscriber`s).
    *   **Chain of Responsibility:** Passed a request along a chain of handlers (expense approval from `Manager` to `Director`).
    *   **Worker Pool:** Managed a large number of tasks with a fixed number of goroutines to control concurrency.
    *   **Fan-in / Fan-out:** Parallelized a pipeline by distributing work (fan-out) and then merging the results (fan-in).
    *   **Rate Limiting:** Controlled the frequency of an operation using a `time.Ticker` as a token dispenser.
*   **Code Location:** `low-level-design/02-design-patterns/`
---
**Session 5: Module 3 - The Building Blocks of Scale**

*   **Summary:** We shifted focus to High-Level Design (HLD), learning the core components required to build scalable systems. We implemented simulations of these components in Go to understand their internal mechanics.
*   **Key Concepts Covered:**
    *   **Scaling:** Understood Vertical (more power) vs. Horizontal (more machines) scaling. Built a dummy backend service to simulate a scalable workload.
    *   **Load Balancing:** Implemented a Layer 7 Load Balancer using Go's `httputil.ReverseProxy`. It used a Round-Robin algorithm to distribute traffic across multiple backend instances.
    *   **Caching:** Implemented the "Cache-Aside" pattern. We simulated a slow DB and a fast in-memory cache, showing how to read from cache first and only hit the DB on a miss.
    *   **Databases & Sharding:** Explored SQL vs. NoSQL. We implemented a Sharding simulation using Consistent Hashing (CRC32) to distribute data across multiple "shard" instances based on the key.
    *   **Message Queues:** Decoupled producers and consumers using Go channels. We built a system where a producer could burst 1000s of requests into a queue, and consumers (workers) processed them at a steady pace without crashing the system.
    *   **Event-Driven Architecture:** Implemented an Event Bus to decouple services. A `UserService` published `UserCreated` events, which triggered independent `Email`, `Analytics`, and `Marketing` services.
    *   **CDNs:** Simulated a global Content Delivery Network. An "Edge Server" cached content from a "Origin Server" to reduce latency for repeated requests.
    *   **Proxies:** Built a Forward Proxy simulation that fetched URLs on behalf of a client, masking the client's identity from the target server.
*   **Code Location:** `high-level-design/03-building-blocks-of-scale/`
---
**Session 6: Module 4 - Reliability & Consistency**

*   **Summary:** We tackled the theoretical limits and practical patterns for keeping distributed systems correct and available. We simulated these distributed system challenges within a single Go process.
*   **Key Concepts Covered:**
    *   **CAP Theorem:** Demonstrated that in the face of a Network Partition (P), you must choose between Consistency (CP - refusing writes) or Availability (AP - accepting writes that may diverge).
    *   **Consistency Patterns:** Simulated "Eventual Consistency". A Master node accepted a write, and we observed the "replication lag" before the Read Replica caught up.
    *   **Availability:** Implemented an "Active-Passive Failover" system. A Load Balancer health-checked the Active server and automatically promoted the Passive server when the Active one "crashed".
    *   **API Design:** Compared REST (full resource, over-fetching), GraphQL (client-specified fields), and gRPC (strict contract) by implementing simulated handlers for a `User` service.
*   **Code Location:** `high-level-design/04-reliability/`
---
**Session 7: Module 5 - Advanced Distributed Concepts**

*   **Summary:** We bridged the gap between intermediate and expert system design by implementing the core algorithms that power modern distributed databases and services.
*   **Key Concepts Covered:**
    *   **Consistent Hashing:** Implemented a Hash Ring with Virtual Nodes to ensure even data distribution and minimal reshuffling when nodes are added or removed.
    *   **Distributed ID Generation:** Built a "Snowflake" ID generator that creates unique, time-sortable 64-bit IDs without a central database, using bitwise manipulation.
    *   **Distributed Transactions:** Implemented the **Saga Pattern** with compensation logic to handle long-running transactions across microservices without locking resources (unlike 2PC).
    *   **Bloom Filters:** Created a probabilistic data structure using double-hashing and bitwise operations to efficiently check for set membership with zero disk I/O.
    *   **Consensus (Raft):** Simulated the Raft Leader Election process, demonstrating how a cluster agrees on a leader and handles network partitions and node failures.
*   **Code Location:** `high-level-design/05-advanced-concepts/`
---
