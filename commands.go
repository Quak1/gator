package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registered map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registered[cmd.Name]
	if !ok {
		return fmt.Errorf("Command not found")
	}

	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registered[name] = f
}
