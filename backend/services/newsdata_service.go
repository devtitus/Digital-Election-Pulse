package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type NewsDataResponse struct {
	Status       string           `json:"status"`
	TotalResults int              `json:"totalResults"`
	Results      []NewsDataResult `json:"results"`
}

type NewsDataResult struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	PubDate     string `json:"pubDate"` // standard pubdate
	SourceID    string `json:"source_id"`
	Description string `json:"description"`
}

func FetchNewsData(ctx context.Context, query string) ([]NewsItem, error) {
	apiKey := os.Getenv("NEWSDATA_API_KEY")
	if apiKey == "" {
		// Fallback or just return empty if not configured (optional source)
		fmt.Println("NEWSDATA_API_KEY not set, skipping NewsData.io fetch")
		return nil, nil
	}
	fmt.Printf("Fetching NewsData.io for query: %s\n", query)

	// Build URL
	baseURL := "https://newsdata.io/api/1/news"
	params := url.Values{}
	params.Add("apikey", apiKey)
	params.Add("q", query)
	params.Add("language", "ta,en") // Tamil and English
	params.Add("country", "in")     // India context

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("newsdata api error: %s", resp.Status)
	}

	var data NewsDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var items []NewsItem
	for _, res := range data.Results {
		// parser time
		// Example: "2023-10-27 10:00:00" or similar? Need to check docs or allow permissive parsing.
		// NewsData often sends "2021-01-29 16:34:02" standard SQL format or RFC.
		// Let's try flexible parsing or default to Now.

		// Attempt parsing standard layouts
		pubDate := time.Now()
		// Try a few formats if needed, or ioTime helper if adapted.

		items = append(items, NewsItem{
			Title:       res.Title,
			Link:        res.Link,
			PublishedAt: pubDate, // simplified for now
			Source:      "NewsData_" + res.SourceID,
		})
	}

	fmt.Printf("Fetched %d items from NewsData.io\n", len(items))
	return items, nil
}
