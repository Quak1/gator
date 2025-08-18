package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Quak1/gator/internal/config"
	"github.com/Quak1/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("error connection to db: %s", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	currentState := state{
		cfg: &cfg,
		db:  dbQueries,
	}

	commands := commands{
		registered: map[string]func(*state, command) error{},
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handleDeleteAllUsers)
	commands.register("users", handleGetUsers)

	if len(os.Args) < 2 {
		log.Fatal("usage: <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	if err := commands.run(&currentState, command{Name: cmdName, Args: cmdArgs}); err != nil {
		log.Fatal(err)
	}
}
