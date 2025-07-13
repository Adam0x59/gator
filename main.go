package main

import (
	"log"
	"os"
	"working/github.com/adam0x59/gator/internal/cli"
	"working/github.com/adam0x59/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	state := cli.State{Config: &cfg}

	commands := cli.Commands{Commands: make(map[string]cli.HandlerFunc)}
	commands.Register("login", cli.HandlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not enough arguments, please try again")
	}
	command := cli.Command{Name: args[1], Args: args[2:]}
	err = commands.Run(&state, command)
	if err != nil {
		log.Fatalln(err)
	}
}
