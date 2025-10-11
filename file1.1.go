package main

import (
	"fmt"
	"log"
	"os"
)

// UserProfile represents a user in the system
type UserProfile struct {
	ID        int
	Username  string
	Email     string
	Active    bool
	CreatedAt string
}

// ValidateUser checks if a user profile is valid
func ValidateUser(user UserProfile) error {
	if user.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if user.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	return nil
}

// ProcessUsers handles a batch of user profiles
func ProcessUsers(users []UserProfile) {
	count := 0
	for _, user := range users {
		if err := ValidateUser(user); err != nil {
			log.Printf("Invalid user %s: %v", user.Username, err)
			continue
		}
		count++
		fmt.Printf("Processing user: %s (%s)\n", user.Username, user.Email)
	}
	fmt.Printf("Successfully processed %d users\n", count)
}

func main() {
	users := []UserProfile{
		{ID: 1, Username: "alice", Email: "alice@example.com", Active: true},
		{ID: 2, Username: "bob", Email: "bob@example.com", Active: false},
		{ID: 3, Username: "charlie", Email: "charlie@example.com", Active: true},
	}

	ProcessUsers(users)
	os.Exit(0)
}
