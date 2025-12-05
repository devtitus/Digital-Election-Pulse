package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type YouTubeComment struct {
	Text        string    `json:"textDisplay"`
	Author      string    `json:"authorDisplayName"`
	PublishedAt time.Time `json:"publishedAt"`
	Likes       int       `json:"likeCount"`
}

type searchResponse struct {
	Items []struct {
		Id struct {
			VideoId string `json:"videoId"`
		} `json:"id"`
	} `json:"items"`
}

type commentThreadResponse struct {
	Items []struct {
		Snippet struct {
			TopLevelComment struct {
				Snippet struct {
					TextDisplay       string `json:"textDisplay"`
					AuthorDisplayName string `json:"authorDisplayName"`
					PublishedAt       string `json:"publishedAt"`
					LikeCount         int    `json:"likeCount"`
				} `json:"snippet"`
			} `json:"topLevelComment"`
		} `json:"snippet"`
	} `json:"items"`
}

func FetchYouTubeComments(query string) ([]YouTubeComment, error) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY is not set")
	}

	// 1. Search for videos (Part 1 of Data Source)
	// Searching for "Party Name speech" or similar as per plan
	searchQuery := fmt.Sprintf("%s speech", query)
	// Increase maxResults to try multiple videos
	searchURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&type=video&q=%s&key=%s&maxResults=5&order=date",
		url.QueryEscape(searchQuery), apiKey)

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("youtube search failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("youtube search failed with status: %d", resp.StatusCode)
	}

	var searchRes searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchRes); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	if len(searchRes.Items) == 0 {
		fmt.Printf("YouTube Search returned 0 videos for query: %s\n", query)
		return []YouTubeComment{}, nil // No videos found
	}

	// 2. Iterate through videos to find one with comments
	for _, item := range searchRes.Items {
		videoId := item.Id.VideoId

		commentsURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/commentThreads?part=snippet&videoId=%s&key=%s&maxResults=50",
			videoId, apiKey)

		cResp, err := http.Get(commentsURL)
		if err != nil {
			fmt.Printf("Failed to fetch comments for video %s: %v\n", videoId, err)
			continue
		}
		defer cResp.Body.Close()

		if cResp.StatusCode != 200 {
			// Comments likely disabled or API error
			fmt.Printf("YouTube Comments API status %d for video %s (likely disabled)\n", cResp.StatusCode, videoId)
			continue
		}

		var commentsRes commentThreadResponse
		if err := json.NewDecoder(cResp.Body).Decode(&commentsRes); err != nil {
			continue
		}

		if len(commentsRes.Items) > 0 {
			// Found comments! Parse and return.
			var comments []YouTubeComment
			for _, cItem := range commentsRes.Items {
				snippet := cItem.Snippet.TopLevelComment.Snippet
				t, _ := time.Parse(time.RFC3339, snippet.PublishedAt)
				comments = append(comments, YouTubeComment{
					Text:        snippet.TextDisplay,
					Author:      snippet.AuthorDisplayName,
					PublishedAt: t,
					Likes:       snippet.LikeCount,
				})
			}
			fmt.Printf("Found %d comments on video %s\n", len(comments), videoId)
			return comments, nil
		}
	}

	fmt.Println("No comments found on any of the recent videos.")
	return []YouTubeComment{}, nil
}
