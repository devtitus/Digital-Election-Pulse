package services

import (
	"context"
	"fmt"

	"github.com/groovili/gogtrends"
)

type TrendItem struct {
	Title   string `json:"title"`
	Traffic string `json:"traffic"`
	Link    string `json:"link"`
}

func FetchTrends(ctx context.Context) ([]TrendItem, error) {
	// "IN-TN" is the geo code for Tamil Nadu, India
	// "TA" is language code for Tamil, but trends API uses slightly different params usually.
	// Daily trends for India (IN) filtered by region isn't always directly exposed in simple calls,
	// checking library capabilities. often defaults to country level.
	// Let's try fetching daily trends for India (IN) and we might have to filter or just show general IN trends.
	// Ideally we want real-time trends for specific interests.

	// Using Daily Trends for India (IN)
	fmt.Printf("Fetching daily trends for geo: IN, category: TA\n")
	dailyTrends, err := gogtrends.Daily(ctx, "IN", "TA") // TA for Tamil
	if err != nil {
		return nil, fmt.Errorf("failed to fetch daily trends: %w", err)
	}

	var items []TrendItem
	for _, trend := range dailyTrends {
		// trend is *gogtrends.TrendingSearch
		// Assuming structure: Title { Query string }, FormattedTraffic string
		title := trend.Title.Query
		traffic := trend.FormattedTraffic

		// Use the first article link if available, or a generic one
		link := "https://trends.google.com/trends/trendingsearches/daily?geo=IN"
		if len(trend.Articles) > 0 {
			link = trend.Articles[0].URL
		}

		items = append(items, TrendItem{
			Title:   title,
			Traffic: traffic,
			Link:    link,
		})
	}

	// Limit to top 10
	if len(items) > 10 {
		items = items[:10]
	}

	fmt.Printf("Fetched %d trending items\n", len(items))
	return items, nil
}
