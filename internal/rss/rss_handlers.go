package rss

import (
	"context"
	"fmt"
	"os"
	"time"
	"working/github.com/adam0x59/gator/internal/cli"
	"working/github.com/adam0x59/gator/internal/database"

	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
)

func HandlerAgg(s *cli.State, cmd cli.Command) error {
	feedStruct, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feedStruct)
	return nil
}

func AddFeed(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no argument provided: " +
			"addFeed requires name and url:\n    gator addfeed <name> <url>")
	} else if len(cmd.Args) > 0 && len(cmd.Args) < 2 {
		return fmt.Errorf("only one argument specified: " +
			"addFeed requires name AND url:\n    gator addfeed <name> <url>")
	}

	user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user data as : %w", err)
	}

	args := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	i, err := s.Db.AddFeed(context.Background(), args)
	if err != nil {
		return fmt.Errorf("error adding feed: %w", err)
	}

	fmt.Printf("Feed %s was added to the database:\n", cmd.Args[0])
	fmt.Printf("Feed Details:\n  id: %s\n  created_at: %s\n  updated_at: %s\n  name: %s\n  url: %s\n, user: %s\n",
		i.ID, i.CreatedAt, i.UpdatedAt, i.Name, i.Url, s.Config.CurrentUserName)

	return nil
}

func Feeds(s *cli.State, cmd cli.Command) error {

	feeds, err := s.Db.Feeds(context.Background())
	if err != nil {
		return err
	}
	//fmt.Println(feeds)
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Feed Name", "URL", "User Name"})
	for _, feed := range feeds {
		feedData := []string{feed.Name, feed.Url, feed.Uname}
		table.Append(feedData)
	}

	table.Render()
	return nil
}
