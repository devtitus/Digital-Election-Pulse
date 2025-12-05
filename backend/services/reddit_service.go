package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type RedditResponse struct {
	Data struct {
		Children []struct {
			Data struct {
				Title     string  `json:"title"`
				Selftext  string  `json:"selftext"`
				Author    string  `json:"author"`
				Url       string  `json:"url"`
				Ups       int     `json:"ups"`
				Created   float64 `json:"created_utc"`
				Subreddit string  `json:"subreddit"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type RedditPost struct {
	Title     string
	Text      string
	URL       string
	Subreddit string
}

func FetchRedditPosts(query string) ([]RedditPost, error) {
	// Subreddits to search
	subreddits := []string{"TamilNadu", "Chennai", "India"}
	var allPosts []RedditPost
	client := &http.Client{Timeout: 10 * time.Second}

	for _, sub := range subreddits {
		// Construct URL: https://www.reddit.com/r/{subreddit}/search.json?q={query}&restrict_sr=1&sort=new&limit=5
		// Need validation/encoding
		encodedQuery := url.QueryEscape(query)
		url := fmt.Sprintf("https://www.reddit.com/r/%s/search.json?q=%s&restrict_sr=1&sort=new&limit=5", sub, encodedQuery)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Error creating Reddit request for r/%s: %v\n", sub, err)
			continue
		}

		// User-Agent is required by Reddit API
		req.Header.Set("User-Agent", "go:election-pulse:v1.0 (by /u/cortex-ai)")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error fetching from r/%s: %v\n", sub, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Printf("Reddit API error r/%s: %d\n", sub, resp.StatusCode)
			continue
		}

		var redditResp RedditResponse
		if err := json.NewDecoder(resp.Body).Decode(&redditResp); err != nil {
			fmt.Printf("Error decoding Reddit response r/%s: %v\n", sub, err)
			continue
		}

		for _, child := range redditResp.Data.Children {
			post := child.Data
			// Basic filtering so we don't capture empty stuff
			if post.Title != "" {
				allPosts = append(allPosts, RedditPost{
					Title:     post.Title,
					Text:      post.Selftext,
					URL:       post.Url,
					Subreddit: post.Subreddit,
				})
			}
		}

		// Respect rate limiting slightly
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("Fetched %d Reddit posts for '%s'\n", len(allPosts), query)
	return allPosts, nil
}
