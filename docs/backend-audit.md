# Backend Audit: The TN Election Pulse

## Executive Summary

The backend is a Go-based microservice that orchestrates real-time data collection from multiple sources (news, social media) and performs AI-powered sentiment analysis for Tamil Nadu political parties. It uses concurrent processing for efficiency, stores results in PostgreSQL, and serves a REST API consumed by the React frontend.

**Key Technologies:**
- **Language**: Go 1.23+
- **Framework**: Fiber (high-performance HTTP framework)
- **Database**: PostgreSQL with GORM ORM
- **AI**: Google Gemini 2.5 Flash API
- **External APIs**: YouTube Data API, NewsData.io, Google News RSS, Reddit JSON API, Google Trends
- **Concurrency**: Goroutines with sync.WaitGroup for parallel data fetching

## Architecture Overview

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   React Frontend │────│   Go Backend     │────│   PostgreSQL    │
│                 │    │   (Fiber)        │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │   External APIs  │
                       │ - YouTube        │
                       │ - News Sources   │
                       │ - Reddit         │
                       │ - Google Trends  │
                       └──────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │   Google Gemini  │
                       │   AI Analysis    │
                       └──────────────────┘
```

### Component Breakdown
- **Entry Point**: `cmd/main.go` - Initializes server and dependencies
- **Routes**: `handlers/routes.go` - REST API endpoints
- **Services**: Modular services for data fetching and AI processing
- **Models**: `models/models.go` - Database entities
- **Database**: `db/` - Connection and schema management

### Data Flow (Analysis Request)
1. User requests analysis via `/analyze` endpoint
2. Orchestrator launches concurrent goroutines for data sources
3. Raw data aggregated into text corpus
4. Corpus sent to Gemini AI for sentiment analysis
5. Results processed, scored (0-100 scale), and stored in DB
6. JSON response returned to frontend

## Entry Point Analysis (`cmd/main.go`)

The application starts with a simple initialization sequence:

```go
func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Connect to Database
    db.Connect()

    // Initialize Fiber app
    app := fiber.New()

    // Middleware
    app.Use(cors.New())

    // Routes
    handlers.SetupRoutes(app)

    // Start server
    log.Fatal(app.Listen(":3000"))
}
```

**Key Points:**
- Environment variables loaded from `.env` (optional)
- Database connection established before server start
- CORS middleware enabled for frontend communication
- Server runs on port 3000

## Database Layer

### Connection Setup (`db/database.go`)
Uses GORM with PostgreSQL driver, connects via `DATABASE_URL` environment variable:

```go
func Connect() {
    dsn := os.Getenv("DATABASE_URL")
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    err = DB.AutoMigrate(&models.Party{}, &models.SentimentSnapshot{})
}
```

### Schema (`db/schema.sql`)
Creates two main tables with seed data:

```sql
CREATE TABLE parties (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    leader VARCHAR(255),
    color_hex VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE sentiment_snapshots (
    id SERIAL PRIMARY KEY,
    party_id INTEGER REFERENCES parties(id),
    score DOUBLE PRECISION,
    key_issue TEXT,
    source_breakdown JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Seed Data**: Pre-populated with major TN parties (DMK, AIADMK, TVK, BJP, NTK).

## Data Models (`models/models.go`)

### Party Model
```go
type Party struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Name      string         `json:"name"`
    Leader    string         `json:"leader"`
    ColorHex  string         `json:"color_hex"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`  // Soft deletes
}
```

### SentimentSnapshot Model
```go
type SentimentSnapshot struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    PartyID         uint      `json:"party_id"`
    Party           Party     `gorm:"foreignKey:PartyID" json:"-"`  // Foreign key
    Score           float64   `json:"score"`                         // 0-100 scale
    KeyTopics       string    `gorm:"type:jsonb" json:"key_topics"`  // JSON array
    Emotion         string    `json:"emotion"`                       // AI-detected emotion
    SourceBreakdown string    `gorm:"type:jsonb" json:"source_breakdown"`
    CreatedAt       time.Time `json:"created_at"`
}
```

**Notes:**
- Uses JSONB for flexible storage of topics and source data
- Foreign key relationship between snapshots and parties
- Soft deletes enabled for parties

## API Layer (`handlers/routes.go`)

### Endpoints
- `GET /api/v1/parties` - List all parties
- `POST /api/v1/analyze` - Trigger sentiment analysis
- `GET /api/v1/latest` - Get latest snapshot for party
- `GET /api/v1/history/:party_id` - Historical data (stub)
- `GET /api/v1/trends` - Google Trends data

### Core Analysis Handler (`AnalyzeParty`)
The main business logic:

```go
func AnalyzeParty(c *fiber.Ctx) error {
    // 1. Parse request
    var req AnalyzeRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }

    // 2. Fetch data concurrently
    data, err := services.FetchAllData(c.Context(), req.PartyName)

    // 3. Build text corpus
    var corpus string
    for _, n := range data.News { corpus += "- " + n.Title + "\n" }
    for _, c := range data.Comments { corpus += "- " + c.Text + "\n" }
    // ... add Reddit posts

    // 4. AI Analysis
    analysis, err := services.AnalyzeSentiment(c.Context(), corpus)

    // 5. Process and store results
    rawScore := analysis.SentimentScore  // -1.0 to 1.0
    finalScore := 50 + (rawScore * 50)   // Convert to 0-100

    // Save to database
    snapshot := models.SentimentSnapshot{...}
    db.DB.Create(&snapshot)

    return c.JSON(analysis)
}
```

**Score Calculation:**
- AI returns sentiment from -1.0 (negative) to 1.0 (positive)
- Formula: `WinningProbability = 50 + (RawScore * 50)`
- Result: 0-100 scale where 50 is neutral

## Core Services Breakdown

### Orchestrator Service (`services/orchestrator.go`)

Coordinates concurrent data fetching from 4 sources:

```go
func FetchAllData(ctx context.Context, partyName string) (*AggregatedData, error) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var data AggregatedData

    wg.Add(4)

    // Launch goroutines for each service
    go func() { defer wg.Done(); rssItems, _ := FetchNews(ctx, partyName); mu.Lock(); data.News = append(data.News, rssItems...); mu.Unlock() }()
    go func() { defer wg.Done(); apiItems, _ := FetchNewsData(ctx, partyName); mu.Lock(); data.News = append(data.News, apiItems...); mu.Unlock() }()
    go func() { defer wg.Done(); data.Comments, _ = FetchYouTubeComments(partyName) }()
    go func() { defer wg.Done(); data.RedditPosts, _ = FetchRedditPosts(partyName) }()

    wg.Wait()
    return &data, nil
}
```

**Key Features:**
- Uses `sync.WaitGroup` for synchronization
- `sync.Mutex` protects shared data structures
- Aggregates results into `AggregatedData` struct
- Continues execution even if some services fail

### AI Service (`services/ai_service.go`)

Handles sentiment analysis via Google Gemini:

```go
func AnalyzeSentiment(ctx context.Context, textData string) (*AIAnalysisResult, error) {
    // Create Gemini client
    client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    model := client.GenerativeModel("gemini-2.5-flash")
    model.ResponseMIMEType = "application/json"

    // Detailed prompt for TN politics
    prompt := fmt.Sprintf(`You are an expert political analyst...`, textData)

    resp, err := model.GenerateContent(ctx, genai.Text(prompt))
    // Parse JSON response...
}
```

**Prompt Engineering:**
- Specialized for Tamil Nadu political context
- Handles local slang and cultural nuances
- Requests structured JSON output with sentiment, emotion, topics
- Includes bias detection and fact-checking instructions

### News Services

#### RSS News Service (`services/news_service.go`)
Fetches from multiple RSS feeds concurrently:

```go
func FetchNews(ctx context.Context, query string) ([]NewsItem, error) {
    urls := []struct{ URL, Source string }{
        {fmt.Sprintf("https://news.google.com/rss/search?q=%s&hl=ta&gl=IN&ceid=IN:ta", query), "Google News"},
        {"https://feeds.feedburner.com/dinamalar/Front_page_news", "Dinamalar"},
        // ... more sources
    }

    // Concurrent fetching with channels
    resultChan := make(chan result, len(urls))
    for _, u := range urls {
        go func() { items, err := fetchFeedItems(ctx, u.URL, u.Source); resultChan <- result{items, err} }()
    }

    // Aggregate results
    var allItems []NewsItem
    for i := 0; i < len(urls); i++ {
        res := <-resultChan
        if res.items != nil { allItems = append(allItems, res.items...) }
    }
}
```

#### NewsData API Service (`services/newsdata_service.go`)
Uses NewsData.io API as alternative/supplement:

```go
func FetchNewsData(ctx context.Context, query string) ([]NewsItem, error) {
    apiKey := os.Getenv("NEWSDATA_API_KEY")
    baseURL := "https://newsdata.io/api/1/news"
    params := url.Values{}
    params.Add("q", query)
    params.Add("language", "ta,en")
    params.Add("country", "in")

    reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
    // HTTP request and JSON parsing...
}
```

### Social Media Services

#### YouTube Service (`services/youtube_service.go`)
Two-step process: search for videos, then fetch comments:

```go
func FetchYouTubeComments(query string) ([]YouTubeComment, error) {
    // 1. Search for videos
    searchQuery := fmt.Sprintf("%s speech", query)
    searchURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?...&q=%s&key=%s&maxResults=5&order=date", 
        url.QueryEscape(searchQuery), apiKey)

    // Parse search results
    var searchRes searchResponse
    json.NewDecoder(resp.Body).Decode(&searchRes)

    // 2. Try to get comments from first video with comments enabled
    for _, item := range searchRes.Items {
        videoId := item.Id.VideoId
        commentsURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/commentThreads?...&videoId=%s&key=%s&maxResults=50", 
            videoId, apiKey)

        // If comments found, parse and return
        var commentsRes commentThreadResponse
        if len(commentsRes.Items) > 0 {
            // Parse comments...
            return comments, nil
        }
    }
}
```

**Resilience Features:**
- Searches multiple videos if comments disabled
- Returns empty array if no comments found (graceful degradation)

#### Reddit Service (`services/reddit_service.go`)
Scrapes Reddit JSON API from multiple subreddits:

```go
func FetchRedditPosts(query string) ([]RedditPost, error) {
    subreddits := []string{"TamilNadu", "Chennai", "India"}

    for _, sub := range subreddits {
        encodedQuery := url.QueryEscape(query)
        url := fmt.Sprintf("https://www.reddit.com/r/%s/search.json?q=%s&restrict_sr=1&sort=new&limit=5", sub, encodedQuery)

        req, _ := http.NewRequest("GET", url, nil)
        req.Header.Set("User-Agent", "go:election-pulse:v1.0 (by /u/cortex-ai)")

        // Parse JSON response...
        var redditResp RedditResponse
        json.NewDecoder(resp.Body).Decode(&redditResp)

        for _, child := range redditResp.Data.Children {
            allPosts = append(allPosts, RedditPost{...})
        }
    }
}
```

### Trends Service (`services/trends_service.go`)
Fetches Google Trends data using gogtrends library:

```go
func FetchTrends(ctx context.Context) ([]TrendItem, error) {
    dailyTrends, err := gogtrends.Daily(ctx, "IN", "TA")  // IN = India, TA = Tamil

    var items []TrendItem
    for _, trend := range dailyTrends {
        items = append(items, TrendItem{
            Title:   trend.Title.Query,
            Traffic: trend.FormattedTraffic,
            Link:    trend.Articles[0].URL,  // First article link
        })
    }
}
```

## Configuration and Dependencies

### Environment Variables
- `DATABASE_URL`: PostgreSQL connection string
- `GEMINI_API_KEY`: Google AI Studio API key
- `YOUTUBE_API_KEY`: Google Cloud Console API key
- `NEWSDATA_API_KEY`: NewsData.io API key (optional)

### Go Module Dependencies (`go.mod`)
```
require (
    github.com/gofiber/fiber/v2 v2.52.5
    github.com/google/generative-ai-go v0.18.0
    gorm.io/gorm v1.25.10
    gorm.io/driver/postgres v1.5.7
    github.com/mmcdole/gofeed v1.3.0
    github.com/groovili/gogtrends v0.0.0-20231215114726-1f8e8a3a3a7f
    github.com/joho/godotenv v1.5.1
)
```

## Error Handling and Resilience

### API Failure Handling
- Services continue execution if individual APIs fail
- Logged errors don't stop the overall process
- Graceful degradation (e.g., empty arrays returned)

### Timeout Management
- HTTP clients have 10-second timeouts
- Context propagation for cancellation

### Rate Limiting Awareness
- Reddit: 500ms delay between subreddit requests
- YouTube: 10,000 units/day quota
- No explicit rate limiting implementation (relies on external limits)

## Performance Considerations

### Concurrency Benefits
- Parallel data fetching reduces total response time
- CPU-bound AI analysis remains sequential

### Resource Usage
- Multiple HTTP connections during data collection
- Memory usage for aggregating large text corpora
- Database connections via GORM connection pool

### Caching Strategy
- Database stores historical snapshots
- API checks for recent data before triggering new analysis
- No in-memory caching implemented

## Security and Best Practices

### API Key Management
- Stored in environment variables
- Not logged or exposed in responses

### Input Validation
- Basic JSON parsing validation
- Query parameter validation for party names

### CORS Configuration
- Enabled for frontend communication
- No additional security headers

## Identified Issues and Recommendations

### Schema vs. Model Mismatch
**Issue**: Database schema has `key_issue TEXT`, but models use `KeyTopics string` (JSONB) and `Emotion string`.

**Impact**: Data storage inconsistency.

**Recommendation**: Update schema.sql to match models:
```sql
ALTER TABLE sentiment_snapshots 
ADD COLUMN key_topics JSONB,
ADD COLUMN emotion VARCHAR(255);
```

### Hardcoded Limits
**Issue**: Fixed limits (5 videos, 50 comments, 5 Reddit posts) scattered in code.

**Recommendation**: Move to configuration or environment variables.

### Error Propagation
**Issue**: Errors logged but not always propagated up properly.

**Recommendation**: Implement structured error handling with error types.

### Testing and Monitoring
**Issue**: No unit tests or monitoring visible.

**Recommendation**: Add Go testing, logging middleware, and metrics collection.

### AI Response Validation
**Issue**: Assumes AI always returns valid JSON.

**Recommendation**: Add schema validation for AI responses.

## Conclusion

The backend demonstrates solid Go fundamentals with effective use of concurrency for performance. The modular service architecture allows for easy extension and maintenance. Key strengths include comprehensive error handling, external API integration, and AI-powered analysis. Areas for improvement focus on consistency, configuration management, and observability.

This audit provides a complete reference for understanding the system's internals and guiding future development.