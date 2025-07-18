package main

import (
	"database/sql"
	"log"
	"os"
	"working/github.com/adam0x59/gator/internal/cli"
	"working/github.com/adam0x59/gator/internal/config"
	"working/github.com/adam0x59/gator/internal/database"
	"working/github.com/adam0x59/gator/internal/rss"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	state := cli.State{Config: &cfg}
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	state.Db = dbQueries
	commands := cli.Commands{Commands: make(map[string]cli.HandlerFunc)}
	commands.Register("login", cli.HandlerLogin)
	commands.Register("register", cli.HandlerRegister)
	commands.Register("reset", cli.HandlerReset)
	commands.Register("users", cli.HandlerGetUsers)
	commands.Register("agg", rss.HandlerAgg)
	commands.Register("addfeed", cli.MiddlewareLoggedIn(rss.AddFeed))
	commands.Register("feeds", rss.Feeds)
	commands.Register("follow", cli.MiddlewareLoggedIn(rss.Follow))
	commands.Register("following", cli.MiddlewareLoggedIn(rss.Following))
	commands.Register("unfollow", cli.MiddlewareLoggedIn(rss.Unfollow))
	commands.Register("browse", cli.MiddlewareLoggedIn(rss.Browse))
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
