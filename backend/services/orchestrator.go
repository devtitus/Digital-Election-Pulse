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
	var data AggregatedData
	var newsErr, youtubeErr, redditErr error

	wg.Add(3)

	// Fetch News
	go func() {
		defer wg.Done()
		data.News, newsErr = FetchNews(ctx, partyName)
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

	return &data, nil
}
