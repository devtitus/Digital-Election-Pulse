package services

import (
	"context"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

type NewsItem struct {
	Title       string
	Link        string
	PublishedAt time.Time
	Source      string
}

func FetchNews(ctx context.Context, query string) ([]NewsItem, error) {
	fp := gofeed.NewParser()
	// Google News RSS URL for Tamil Nadu/Tamil language context
	rssURL := fmt.Sprintf("https://news.google.com/rss/search?q=%s&hl=ta&gl=IN&ceid=IN:ta", query)

	// Create a derived context with timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	feed, err := fp.ParseURLWithContext(rssURL, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed: %w", err)
	}

	var items []NewsItem
	for _, item := range feed.Items {
		// Limit to top 10 as per plan
		if len(items) >= 10 {
			break
		}

		pubDate := ioTime(item.PublishedParsed)

		items = append(items, NewsItem{
			Title:       item.Title,
			Link:        item.Link,
			PublishedAt: pubDate,
			Source:      "Google News",
		})
	}
	return items, nil
}

func ioTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Now()
}
