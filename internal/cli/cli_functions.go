package cli

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"working/github.com/adam0x59/gator/internal/database"

	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no argument provided: " +
			"login requires a username. command format:\n    gator login <username>")
	}
	_, err := s.Db.GetUser(context.Background(), sql.NullString{String: cmd.Args[0], Valid: true})
	if err == sql.ErrNoRows {
		return fmt.Errorf("user %s does not exist: error as %w", cmd.Args[0], err)
	} else if err != nil {
		return err
	}
	s.Config.SetUser(cmd.Args[0])
	fmt.Printf("User has been set to %s\n", s.Config.CurrentUserName)
	fmt.Printf("\nCurrent Config:\n")
	fmt.Printf("Database URL: %s\n", s.Config.DbUrl)
	fmt.Printf("User Name: %s\n", s.Config.CurrentUserName)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no argument provided: " +
			"register requires name of user to be registered:\n    gator register <name>")
	}

	id := uuid.NullUUID{UUID: uuid.New(), Valid: true}
	created_at := sql.NullTime{Time: time.Now(), Valid: true}
	updated_at := created_at
	name := sql.NullString{String: cmd.Args[0], Valid: true}

	args := database.CreateUserParams{
		Column1: id,
		Column2: created_at,
		Column3: updated_at,
		Column4: name,
	}
	_, err := s.Db.GetUser(context.Background(), name)
	if err == sql.ErrNoRows {
		userData, err := s.Db.CreateUser(context.Background(), args)
		if err != nil {
			return err
		}
		fmt.Printf("User %s was added to the database\n", cmd.Args[0])
		fmt.Printf("User detail:\n  id: %s\n  created_at: %s\n  updated_at: %s\n  name: %s\n",
			userData.ID, userData.CreatedAt.Time, userData.UpdatedAt.Time, userData.Name)
		loginErr := HandlerLogin(s, cmd)
		if loginErr != nil {
			return loginErr
		}
	} else if err != nil {
		return err
	} else {
		return fmt.Errorf("user %s already exists in the database", cmd.Args[0])
	}
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	err := s.Db.Reset(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Database reset successful!")
	return nil
}

func HandlerGetUsers(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return fmt.Errorf("no users found in the database")
	}
	for _, user := range users {
		if user == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}
