package services

import (
	"context"
	"fmt"
	"net/http"
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
	// RSS Sources
	// 1. Google News (Tamil Nadu context)
	// 2. Dinamalar (Front Page)
	// 3. Dinakaran (Latest)

	urls := []struct {
		URL    string
		Source string
	}{
		{
			URL:    fmt.Sprintf("https://news.google.com/rss/search?q=%s&hl=ta&gl=IN&ceid=IN:ta", query),
			Source: "Google News",
		},
		{
			URL:    "https://feeds.feedburner.com/dinamalar/Front_page_news",
			Source: "Dinamalar",
		},
		{
			URL:    "https://tamil.hindustantimes.com/rss/tamilnadu",
			Source: "Hindustan Times Tamil",
		},
		{
			URL:    "https://tamil.news18.com/commonfeeds/v1/tam/rss/tamil-nadu.xml",
			Source: "News18 Tamil Nadu",
		},
	}

	// Channel to collect results
	type result struct {
		items []NewsItem
		err   error
	}
	resultChan := make(chan result, len(urls))

	for _, u := range urls {
		go func(urlInfo struct{ URL, Source string }) {
			// Fetch with explicit headers
			items, err := fetchFeedItems(ctx, urlInfo.URL, urlInfo.Source)
			if err != nil {
				// Log error but don't fail everything
				fmt.Printf("Error fetching %s: %v\n", urlInfo.Source, err)
				resultChan <- result{nil, err}
				return
			}
			resultChan <- result{items, nil}
		}(u)
	}

	var allItems []NewsItem
	for i := 0; i < len(urls); i++ {
		res := <-resultChan
		if res.items != nil {
			allItems = append(allItems, res.items...)
		}
	}

	return allItems, nil
}

func fetchFeedItems(ctx context.Context, url, source string) ([]NewsItem, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers to mimic a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "https://www.google.com/")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http error: %s", resp.Status)
	}

	fp := gofeed.NewParser()
	feed, err := fp.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var items []NewsItem
	for _, item := range feed.Items {
		pubDate := ioTime(item.PublishedParsed)
		items = append(items, NewsItem{
			Title:       item.Title,
			Link:        item.Link,
			PublishedAt: pubDate,
			Source:      source,
		})
		// Limit per source
		if len(items) >= 5 {
			break
		}
	}
	return items, nil
}

func ioTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Now()
}
