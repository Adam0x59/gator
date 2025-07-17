package rss

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	"working/github.com/adam0x59/gator/internal/cli"
	"working/github.com/adam0x59/gator/internal/database"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

func AddFeed(s *cli.State, cmd cli.Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no argument provided: " +
			"addFeed requires name and url:\n    gator addfeed <name> <url>")
	} else if len(cmd.Args) > 0 && len(cmd.Args) < 2 {
		return fmt.Errorf("only one argument specified: " +
			"addFeed requires name AND url:\n    gator addfeed <name> <url>")
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

	followCmd := cli.Command{Name: "follow", Args: []string{i.Url}}
	err = Follow(s, followCmd, user)
	if err != nil {
		return err
	}

	return nil
}

func Feeds(s *cli.State, cmd cli.Command) error {

	feeds, err := s.Db.Feeds(context.Background())
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Feed Name", "URL", "User Name"})
	for _, feed := range feeds {
		feedData := []string{feed.Name, feed.Url, feed.Uname}
		table.Append(feedData)
	}

	table.Render()
	return nil
}

func Follow(s *cli.State, cmd cli.Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no argument provided: " +
			"follow requires a url:\n    gator follow <url>")
	}
	url := cmd.Args[0]
	feedID, err := s.Db.Feed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed id as: %w", err)
	}

	args := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedID,
	}
	var pqErr *pq.Error
	i, err := s.Db.CreateFeedFollow(context.Background(), args)
	if err != nil {
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			fmt.Println("You already follow this feed!")
			return nil
		}
		return fmt.Errorf("error creating feed-follow as: %w", err)
	}
	fmt.Printf("User \"%s\", is now following \"%s\"\n", i.UserName, i.FeedName)
	return nil
}

func Following(s *cli.State, cmd cli.Command, user database.User) error {
	feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting feeds for user %s: %w", s.Config.CurrentUserName, err)
	}
	fmt.Printf("User %q follows:\n", s.Config.CurrentUserName)
	for _, feed := range feeds {
		fmt.Printf(" - %s\n", feed.FeedName)
	}
	return nil
}

func Unfollow(s *cli.State, cmd cli.Command, user database.User) error {
	feed, err := s.Db.Feed(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error getting feed id: %w", err)
	}
	args := database.DeleteFollowParams{
		UserID: user.ID,
		FeedID: feed,
	}

	err = s.Db.DeleteFollow(context.Background(), args)
	if err != nil {
		fmt.Errorf("error deleting feed %s, for %q: %w", cmd.Args[0], user.Name, err)
	}
	fmt.Printf("%q is no longer following %q\n", user.Name, cmd.Args[0])
	return nil
}
