package main

import (
	"fmt"

	"github.com/Quak1/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	cfg.SetUser("quak")

	cfg, err = config.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
