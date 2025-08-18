package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Quak1/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	currentState := state{
		cfg: &cfg,
	}

	commands := commands{
		registered: map[string]func(*state, command) error{},
	}
	commands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Command name required")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	if err := commands.run(&currentState, command{Name: cmdName, Args: cmdArgs}); err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)
}
