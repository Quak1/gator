package main

import (
	"context"

	"github.com/Quak1/gator/internal/database"
)

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
		if err != nil {
			return err
		}

		return handler(s, c, user)
	}
}
