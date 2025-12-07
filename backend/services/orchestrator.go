package services

import (
	"context"
	"fmt"
	"sync"
)

type AggregatedData struct {
	News        []NewsItem
	Comments    []YouTubeComment
	RedditPosts []RedditPost
}

func FetchAllData(ctx context.Context, partyName string) (*AggregatedData, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex for shared data.News access
	var data AggregatedData
	var newsErr, youtubeErr, redditErr, newsDataErr error

	wg.Add(4)

	// Fetch News (RSS)
	go func() {
		defer wg.Done()
		fmt.Println("Starting RSS News fetch...")
		rssItems, err := FetchNews(ctx, partyName)
		newsErr = err

		mu.Lock()
		if rssItems != nil {
			data.News = append(data.News, rssItems...)
		}
		mu.Unlock()
	}()

	// Fetch NewsData.io (API)
	go func() {
		defer wg.Done()
		fmt.Println("Starting NewsData.io fetch...")
		apiItems, err := FetchNewsData(ctx, partyName)
		newsDataErr = err

		mu.Lock()
		if apiItems != nil {
			data.News = append(data.News, apiItems...)
		}
		mu.Unlock()
	}()

	// Fetch YouTube
	go func() {
		defer wg.Done()
		data.Comments, youtubeErr = FetchYouTubeComments(partyName)
	}()

	// Fetch Reddit
	go func() {
		defer wg.Done()
		data.RedditPosts, redditErr = FetchRedditPosts(partyName)
	}()

	wg.Wait()

	if newsErr != nil {
		fmt.Printf("News fetch error: %v\n", newsErr)
	}
	if youtubeErr != nil {
		fmt.Printf("YouTube fetch error: %v\n", youtubeErr)
	}
	if redditErr != nil {
		fmt.Printf("Reddit fetch error: %v\n", redditErr)
	}
	if newsDataErr != nil {
		fmt.Printf("NewsData fetch error: %v\n", newsDataErr)
	}

	return &data, nil
}
