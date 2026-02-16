package main

import (
	"errors"
	"fmt"
)

// --- Strategy 2: Custom Error Types ---
// By creating our own error type, we can add more structured information
// to our errors. This is much better than just returning simple error strings.
type UserNotFoundError struct {
	UserID string
}

// The Error() method makes our struct satisfy the built-in `error` interface.
func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user with ID '%s' not found", e.UserID)
}

// A mock database of users.
var userDatabase = map[string]string{
	"123": "Alice",
	"456": "Bob",
}

// FindUser simulates looking for a user in a database.
// It can return our custom error type.
func FindUser(userID string) (string, error) {
	if user, ok := userDatabase[userID]; ok {
		return user, nil
	}
	// Return an instance of our custom error.
	return "", &UserNotFoundError{UserID: userID}
}

// GetUserProfile simulates a higher-level operation that uses FindUser.
func GetUserProfile(userID string) (string, error) {
	username, err := FindUser(userID)
	if err != nil {
		// --- Strategy 3: Error Wrapping ---
		// Instead of just returning the original error, we wrap it with more
		// context. The `%w` verb in fmt.Errorf is crucial; it creates an
		// error chain. This lets us inspect the underlying error later.
		return "", fmt.Errorf("could not get user profile: %w", err)
	}
	return fmt.Sprintf("Profile for %s", username), nil
}

func main() {
	// --- Scenario 1: User Found ---
	fmt.Println("Searching for user 123...")
	profile, err := GetUserProfile("123")
	// --- Strategy 1: Simple Error Checking ---
	// This is the most common pattern in Go. Check if err is not nil.
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: %s\n", profile)
	}
	fmt.Println("---")

	// --- Scenario 2: User Not Found ---
	fmt.Println("Searching for user 999...")
	profile, err = GetUserProfile("999")
	if err != nil {
		fmt.Printf("Error occurred: %v\n", err)

		// Now we can inspect the error chain to see *why* it failed.
		// `errors.As` checks if any error in the chain matches our custom type.
		var userNotFoundErr *UserNotFoundError
		if errors.As(err, &userNotFoundErr) {
			fmt.Println("Error Type: This was a 'User Not Found' error.")
			fmt.Printf("Specifics: Could not find user with ID: %s\n", userNotFoundErr.UserID)
		}
	} else {
		fmt.Printf("Success: %s\n", profile)
	}
}
