package cli

import "fmt"

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no argument provided: " +
			"login requires a username. command format:\n    gator login <username>")
	}
	s.Config.SetUser(cmd.Args[0])
	fmt.Printf("User has been set to %s\n", s.Config.CurrentUserName)
	fmt.Printf("\nCurrent Config:\n")
	fmt.Printf("Database URL: %s\n", s.Config.DbUrl)
	fmt.Printf("User Name: %s\n", s.Config.CurrentUserName)
	return nil
}
