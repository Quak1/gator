package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/Quak1/gator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "quak-gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting rss feed: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading rss feed: %w", err)
	}

	feed := &RSSFeed{}
	if err := xml.Unmarshal(data, feed); err != nil {
		return nil, fmt.Errorf("error parsing rss feed: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}

	return feed, nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:            nextFeed.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:     time.Now(),
	})
	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {
		pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return fmt.Errorf("error parsing time for post %s: %w", item.Title, err)
		}

		err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			FeedID:      nextFeed.ID,
			PublishedAt: pubDate,
		})
		if err != nil {
			fmt.Printf("error saving post: %v\n", err)
		}
	}

	fmt.Printf("processed %d posts from '%s'\n", len(feed.Channel.Item), feed.Channel.Title)

	return nil
}
