package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// RaftState represents the node's current role.
type RaftState int

const (
	Follower RaftState = iota
	Candidate
	Leader
)

func (s RaftState) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	default:
		return "Unknown"
	}
}

// RequestVoteArgs mimics the RPC arguments.
type RequestVoteArgs struct {
	Term        int
	CandidateID int
}

// RequestVoteReply mimics the RPC reply.
type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

// HeartbeatArgs mimics the AppendEntries RPC.
type HeartbeatArgs struct {
	Term     int
	LeaderID int
}

// RaftNode represents a single server in the cluster.
type RaftNode struct {
	mu sync.Mutex
	id int

	// Persistent state
	currentTerm int
	votedFor    int // -1 means null

	// Volatile state
	state        RaftState
	electionTimer *time.Timer
	peers        []*RaftNode // Network simulation: pointers to other nodes

	// Network Channels (Simulating RPCs)
	voteReqChan   chan RequestVoteArgs
	voteReplyChan chan RequestVoteReply
	heartbeatChan chan HeartbeatArgs
}

func NewRaftNode(id int) *RaftNode {
	rn := &RaftNode{
		id:            id,
		state:         Follower,
		votedFor:      -1,
		voteReqChan:   make(chan RequestVoteArgs, 10),
		voteReplyChan: make(chan RequestVoteReply, 10),
		heartbeatChan: make(chan HeartbeatArgs, 10),
	}
	return rn
}

// Start begins the main loop of the Raft node.
func (rn *RaftNode) Start() {
	go func() {
		// Initial random timeout to prevent split votes
		rn.resetElectionTimer()

		for {
			select {
			case <-rn.electionTimer.C:
				// Timeout! Time to start an election.
				rn.startElection()

			case args := <-rn.voteReqChan:
				rn.handleRequestVote(args)

			case args := <-rn.heartbeatChan:
				rn.handleHeartbeat(args)
			}
		}
	}()
}

func (rn *RaftNode) resetElectionTimer() {
	rn.mu.Lock()
	defer rn.mu.Unlock()
	
	if rn.electionTimer != nil {
		rn.electionTimer.Stop()
	}
	// Random timeout between 150ms and 300ms
	duration := time.Duration(150+rand.Intn(150)) * time.Millisecond
	rn.electionTimer = time.NewTimer(duration)
}

// --- ELECTION LOGIC ---

func (rn *RaftNode) startElection() {
	rn.mu.Lock()
	if rn.state == Leader {
		rn.mu.Unlock()
		return
	}
	
	// Transition to Candidate
	rn.state = Candidate
	rn.currentTerm++
	rn.votedFor = rn.id
	term := rn.currentTerm
	rn.mu.Unlock()

	fmt.Printf("[%d] Election Timeout! Becoming Candidate (Term %d)\n", rn.id, term)

	
	for _, peer := range rn.peers {
		if peer.id == rn.id {
			continue
		}
		
		// Simulate sending RPC
		go func(p *RaftNode) {
			p.voteReqChan <- RequestVoteArgs{Term: term, CandidateID: rn.id}
		}(peer)
	}

	// In a real implementation, we would wait for replies here.
	// For this simulation, we'll assume we get votes if we are the first candidate.
	// We'll simulate the "Vote Counting" logic simply by checking back later.
}

func (rn *RaftNode) handleRequestVote(args RequestVoteArgs) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	reply := RequestVoteReply{Term: rn.currentTerm, VoteGranted: false}

	if args.Term > rn.currentTerm {
		fmt.Printf("[%d] Saw higher term %d. Stepping down to Follower.\n", rn.id, args.Term)
		rn.currentTerm = args.Term
		rn.state = Follower
		rn.votedFor = -1
	}

	if args.Term < rn.currentTerm {
		// Reject outdated term
		return 
	}

	if rn.votedFor == -1 || rn.votedFor == args.CandidateID {
		rn.votedFor = args.CandidateID
		reply.VoteGranted = true
		fmt.Printf("[%d] Voted for Candidate %d (Term %d)\n", rn.id, args.CandidateID, args.Term)
		
		// Reset timer because we heard from a valid candidate
		if rn.electionTimer != nil {
			rn.electionTimer.Reset(time.Duration(150+rand.Intn(150)) * time.Millisecond)
		}
	}
	
	// In a full impl, we'd send 'reply' back. 
	// Here, we'll simulate the "Winning" condition:
	// If we vote for someone, and they get enough votes, they become leader.
	if reply.VoteGranted {
		// Magic shortcut for simulation: Tell the candidate they got a vote
		go func() {
			// (Simplification: Accessing peer state directly to count votes)
			// In reality: network reply -> count -> check quorum
			candidate := rn.getPeer(args.CandidateID)
			candidate.receiveVote()
		}()
	}
}

func (rn *RaftNode) receiveVote() {
	rn.mu.Lock()
	defer rn.mu.Unlock()
	
	if rn.state != Candidate {
		return
	}
	
	// We assume we got a vote (simplified).
	// Real implementation tracks count.
	// We will simulate "Winning" if we just assume we got majority for this demo.
	// Let's say getting *one* external vote + self vote = 2/3 = Majority.
	
	fmt.Printf("[%d] Received Vote! Achieving Quorum. Becoming LEADER.\n", rn.id)
	rn.state = Leader
	
	// Start Heartbeat loop
	go rn.sendHeartbeats()
}

// --- LEADER LOGIC ---

func (rn *RaftNode) sendHeartbeats() {
	for {
		rn.mu.Lock()
		if rn.state != Leader {
			rn.mu.Unlock()
			return
		}
		term := rn.currentTerm
		rn.mu.Unlock()

		// fmt.Printf("[%d] Leader sending heartbeats...\n", rn.id)

		for _, peer := range rn.peers {
			if peer.id == rn.id { continue }
			peer.heartbeatChan <- HeartbeatArgs{Term: term, LeaderID: rn.id}
		}
		
		time.Sleep(50 * time.Millisecond) // Heartbeat interval
	}
}

func (rn *RaftNode) handleHeartbeat(args HeartbeatArgs) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	if args.Term >= rn.currentTerm {
		if rn.state != Follower {
			fmt.Printf("[%d] Discovered Leader %d (Term %d). Stepping down.\n", rn.id, args.LeaderID, args.Term)
		}
		rn.currentTerm = args.Term
		rn.state = Follower
		
		// Reset timer because Leader is alive
		if rn.electionTimer != nil {
			rn.electionTimer.Reset(time.Duration(150+rand.Intn(150)) * time.Millisecond)
		}
	}
}

// Helper to find peer object
func (rn *RaftNode) getPeer(id int) *RaftNode {
	for _, p := range rn.peers {
		if p.id == id { return p }
	}
	return nil
}

func main() {
	// Create 3 nodes
	nodes := []*RaftNode{
		NewRaftNode(1),
		NewRaftNode(2),
		NewRaftNode(3),
	}

	// Connect them (Mesh Topology)
	for _, n := range nodes {
		n.peers = nodes
	}

	// Start them
	fmt.Println("--- Starting Raft Cluster (3 Nodes) ---")
	for _, n := range nodes {
		n.Start()
	}

	// Wait for Leader Election
	time.Sleep(2 * time.Second)

	// Simulate Partition: Kill the Leader
	var leader *RaftNode
	for _, n := range nodes {
		n.mu.Lock()
		if n.state == Leader {
			leader = n
		}
		n.mu.Unlock()
	}

	if leader != nil {
		fmt.Printf("\n--- CRASHING LEADER Node %d ---\n", leader.id)
		// Stop sending heartbeats (Simulate crash)
		leader.mu.Lock()
		leader.state = Follower // Just demote to stop the heartbeat loop
		leader.electionTimer.Stop() // Stop it from re-electing
		leader.mu.Unlock()
	}

	// Wait for Re-Election
	fmt.Println("Waiting for new election...")
	time.Sleep(2 * time.Second)
	
	fmt.Println("\n--- Simulation End ---")
}
