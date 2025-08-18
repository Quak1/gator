package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Quak1/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <feed_url>", cmd.Name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("error creating feed, user not found: %w", err)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:      name,
		Url:       url,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}

	for _, feed := range feeds {
		printFeed(feed)
	}

	return nil
}

func printFeed(feed database.GetFeedsRow) {
	fmt.Printf(" - ID:        %s\n", feed.ID)
	fmt.Printf(" - CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf(" - UpdatedAt: %v\n", feed.UpdatedAt)
	fmt.Printf(" - Name:      %s\n", feed.Name)
	fmt.Printf(" - URL:       %s\n", feed.Url)
	fmt.Printf(" - User:      %s\n", feed.Username)
	fmt.Printf(" - UserID:    %s\n", feed.UserID)
	fmt.Println()
}
