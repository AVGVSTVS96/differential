package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// UserProfile represents a user in the system with enhanced fields
type UserProfile struct {
	ID        int
	Username  string
	Email     string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      string
}

// ValidateUser checks if a user profile is valid with additional rules
func ValidateUser(user UserProfile) error {
	if user.Username == "" || len(user.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	if user.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if user.Role == "" {
		return fmt.Errorf("role must be specified")
	}
	return nil
}

// ProcessUsers handles a batch of user profiles with improved logging
func ProcessUsers(users []UserProfile) int {
	count := 0
	failed := 0
	for _, user := range users {
		if err := ValidateUser(user); err != nil {
			log.Printf("Validation failed for user %s: %v", user.Username, err)
			failed++
			continue
		}
		count++
		fmt.Printf("Processing user: %s (%s) - Role: %s\n", user.Username, user.Email, user.Role)
	}
	fmt.Printf("Successfully processed %d users, %d failed\n", count, failed)
	return count
}

func main() {
	users := []UserProfile{
		{ID: 1, Username: "alice", Email: "alice@example.com", Active: true, Role: "admin"},
		{ID: 2, Username: "bob", Email: "bob@example.com", Active: false, Role: "user"},
		{ID: 3, Username: "charlie", Email: "charlie@example.com", Active: true, Role: "moderator"},
		{ID: 4, Username: "diana", Email: "diana@example.com", Active: true, Role: "user"},
	}

	processed := ProcessUsers(users)
	log.Printf("Total processed: %d", processed)
	os.Exit(0)
}
