package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// User represents our database model. It has many fields.
type User struct {
	ID        int
	Name      string
	Email     string
	Age       int
	Address   string
	CreatedAt time.Time
}

var dbUser = User{
	ID:        1,
	Name:      "Alice",
	Email:     "alice@example.com",
	Age:       30,
	Address:   "123 Go Lane",
	CreatedAt: time.Now(),
}

// --- 1. REST Simulation ---
// GET /users/1
// Returns the full resource representation (JSON).
func HandleREST() string {
	fmt.Println("[REST] Fetching User resource...")
	// Marshals the entire struct.
	jsonData, _ := json.MarshalIndent(dbUser, "", "  ")
	return string(jsonData)
}

// --- 2. GraphQL Simulation ---
// POST /graphql { query: "user(id:1) { name, email }" }
// Returns only what was asked for.
func HandleGraphQL(requestedFields []string) string {
	fmt.Printf("[GraphQL] Fetching User fields: %v...\n", requestedFields)
	
	result := make(map[string]interface{})
	
	for _, field := range requestedFields {
		switch strings.ToLower(field) {
		case "id":
			result["id"] = dbUser.ID
		case "name":
			result["name"] = dbUser.Name
		case "email":
			result["email"] = dbUser.Email
		case "age":
			result["age"] = dbUser.Age
		// Note: Address and CreatedAt are skipped if not requested!
		}
	}

	jsonData, _ := json.MarshalIndent(result, "", "  ")
	return string(jsonData)
}

// --- 3. gRPC Simulation ---
// rpc GetUser(GetUserRequest) returns (UserResponse);
// Strictly typed request and response structs.
// In reality, this would be binary data (Protobuf), not printed structs.

type GrpcUserRequest struct {
	UserID int32
}

type GrpcUserResponse struct {
	Name  string
	Email string
	// Notice: This specific RPC response might NOT include 'Address' or 'Age'
	// if the proto definition was designed for a lightweight view.
}

func HandleGRPC(req GrpcUserRequest) GrpcUserResponse {
	fmt.Printf("[gRPC] Calling GetUser procedure with ID=%d...\n", req.UserID)
	
	// Logic to map DB model to Proto response
	return GrpcUserResponse{
		Name:  dbUser.Name,
		Email: dbUser.Email,
	}
}

func main() {
	fmt.Println("--- Scenario: Mobile App needs just Name and Email ---")

	fmt.Println("\n--- 1. REST Approach ---")
	// Problem: Over-fetching. We get Age, Address, CreatedAt even though we don't need them.
	// Bandwidth wasted.
	restResponse := HandleREST()
	fmt.Println(restResponse)

	fmt.Println("\n--- 2. GraphQL Approach ---")
	// Solution: We ask exactly for what we need.
	// Bandwidth saved. Flexible for different clients.
	gqlResponse := HandleGraphQL([]string{"name", "email"})
	fmt.Println(gqlResponse)

	fmt.Println("\n--- 3. gRPC Approach ---")
	// Solution: Strongly typed contract.
	// Extremely fast parsing (no JSON overhead in real life).
	// Great for microservice-to-microservice communication.
	grpcResponse := HandleGRPC(GrpcUserRequest{UserID: 1})
	fmt.Printf("Response Struct: %+v\n", grpcResponse)
}
