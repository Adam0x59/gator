package cli

import (
	"fmt"
	"working/github.com/adam0x59/gator/internal/config"
	"working/github.com/adam0x59/gator/internal/database"
)

type State struct {
	Db     *database.Queries
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type HandlerFunc func(*State, Command) error

type Commands struct {
	Commands map[string]HandlerFunc
}

func (c *Commands) Run(s *State, cmd Command) error {
	function, ok := c.Commands[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	err := function(s, cmd)
	if err != nil {
		return fmt.Errorf("command %s failed with error: %w", cmd.Name, err)
	}
	return nil
}

func (c *Commands) Register(name string, f HandlerFunc) {
	c.Commands[name] = f
}
