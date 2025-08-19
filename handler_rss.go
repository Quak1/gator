package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Quak1/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <interval>", cmd.Name)
	}
	interval := cmd.Args[0]

	duration, err := time.ParseDuration(interval)
	if err != nil {
		return fmt.Errorf("error parsing interval: %w", err)
	}

	fmt.Printf("Printing feeds every %s\n", duration)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	scrapeFeeds(s)
	for range ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <feed_url>", cmd.Name)
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

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		Name:      user.Name,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
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

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}
	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed: %w", err)
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		Name:      s.cfg.CurrentUsername,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}

	fmt.Printf("Feed name: %s\n", follow.FeedName)
	fmt.Printf("User name: %s\n", follow.UserName)

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}
	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed: %w", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error unfollowing feed: %w", err)
	}

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	feeds, err := s.db.GetFeedFolllowsForUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("error getting following feeds for user '%s': %w", s.cfg.CurrentUsername, err)
	}

	for _, feed := range feeds {
		fmt.Printf(" - %s\n", feed.FeedName)
	}

	return nil
}
