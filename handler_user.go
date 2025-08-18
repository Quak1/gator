package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Quak1/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	username := cmd.Args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("error getting user '%s': %w", username, err)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("Couldn't set user: %w", err)
	}

	fmt.Printf("Username '%s' has been set\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <name>", cmd.Name)
	}
	username := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		Name:      username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error creating user '%s': %w", username, err)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("Couldn't set user: %w", err)
	}

	fmt.Printf("user '%s' created!\n", user.Name)

	return nil
}

func handleDeleteAllUsers(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("error clearing users table: %w", err)
	}

	fmt.Println("users table cleared successfully")

	return nil
}

func handleGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	for _, user := range users {
		name := user.Name
		fmt.Printf("* %s", name)
		if name == s.cfg.CurrentUsername {
			fmt.Printf(" (current)")
		}
		fmt.Println()
	}

	return nil
}
