package handlers

import (
	"election-pulse-backend/db"
	"election-pulse-backend/models"
	"election-pulse-backend/services"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/parties", GetParties)
	api.Post("/analyze", AnalyzeParty)
	api.Get("/latest", GetLatestSnapshot)
	api.Get("/history/:party_id", GetHistory)
	api.Get("/trends", GetTrends)
}

func GetParties(c *fiber.Ctx) error {
	var parties []models.Party
	if result := db.DB.Find(&parties); result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
	}
	return c.JSON(parties)
}

type AnalyzeRequest struct {
	PartyName string `json:"party_name"`
}

func AnalyzeParty(c *fiber.Ctx) error {
	var req AnalyzeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// 1. Fetch Data
	data, err := services.FetchAllData(c.Context(), req.PartyName)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data: " + err.Error()})
	}

	// 2. Prepare text for AI
	var corpus string
	corpus += "Latest News Headlines:\n"
	for _, n := range data.News {
		corpus += "- " + n.Title + "\n"
	}
	corpus += "\nSocial Media Comments:\n"
	for _, c := range data.Comments {
		corpus += "- " + c.Text + "\n"
	}

	corpus += "\nReddit Discussions:\n"
	for _, p := range data.RedditPosts {
		corpus += fmt.Sprintf("- Title: %s\n  Body: %s\n", p.Title, p.Text)
	}

	fmt.Printf("Corpus prepared: %d news, %d comments, %d reddit posts\n", len(data.News), len(data.Comments), len(data.RedditPosts))

	// 3. Analyze with AI
	analysis, err := services.AnalyzeSentiment(c.Context(), corpus)
	if err != nil {
		fmt.Printf("Error analyzing sentiment: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "AI Analysis failed: " + err.Error()})
	}

	// 4. Save Snapshot (Optional: Find Party ID first)
	var party models.Party
	db.DB.Where("name = ?", req.PartyName).First(&party)
	if party.ID != 0 {
		// Marshal KeyTopics to JSON string
		keyTopicsJSON, _ := json.Marshal(analysis.KeyTopics)

		snapshot := models.SentimentSnapshot{
			PartyID: party.ID,
			Score:   analysis.SentimentScore * 100, // Convert -1.0-1.0 to 0-100 if needed, or mapping logic?
			// The plan said: WinningProbability = 50 + (RawScore * 50)
			// RawScore assumed to be -1 to 1.
			// Let's apply: 50 + (analysis.SentimentScore * 50)
			// Wait, the previous prompt Plan said: 0-100 score.
			// Let's assume the AI returns -1.0 to 1.0.
			KeyTopics:       string(keyTopicsJSON),
			Emotion:         analysis.Emotion,
			SourceBreakdown: "{}", // simplification
			CreatedAt:       time.Now(),
		}
		// Adjust score logic
		rawScore := analysis.SentimentScore
		// Clamp between -1 and 1 just in case
		if rawScore > 1 {
			rawScore = 1
		}
		if rawScore < -1 {
			rawScore = -1
		}

		finalScore := 50 + (rawScore * 50)
		snapshot.Score = finalScore
		// snapshot.KeyIssue = safeGetTopic(analysis.KeyTopics) // Removed

		db.DB.Create(&snapshot)

		// Return result with the calculated score
		analysis.SentimentScore = finalScore
	}

	return c.JSON(analysis)
}

func GetLatestSnapshot(c *fiber.Ctx) error {
	partyName := c.Query("party_name")
	if partyName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Party Name is required"})
	}

	var party models.Party
	if err := db.DB.Where("name = ?", partyName).First(&party).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Party not found"})
	}

	var snapshot models.SentimentSnapshot
	// Get the latest snapshot
	if err := db.DB.Where("party_id = ?", party.ID).Order("created_at desc").First(&snapshot).Error; err != nil {
		// No data found is not an error here, just return empty Analysis
		// Reuse AIAnalysisResult structure for consistency or separate?
		// Let's return a specific structure for UI
		return c.JSON(fiber.Map{"exists": false})
	}

	// Unmarshal KeyTopics
	var keyTopics []string
	if snapshot.KeyTopics != "" {
		json.Unmarshal([]byte(snapshot.KeyTopics), &keyTopics)
	}

	// Map DB snapshot to response format matching AnalyzeParty
	return c.JSON(fiber.Map{
		"exists":          true,
		"sentiment_score": snapshot.Score,
		"key_topics":      keyTopics,
		"emotion":         snapshot.Emotion,
		"created_at":      snapshot.CreatedAt,
	})
}

func safeGetTopic(topics []string) string {
	if len(topics) > 0 {
		return topics[0]
	}
	return "General"
}

func GetHistory(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get History - To Be Implemented"})
}

func GetTrends(c *fiber.Ctx) error {
	trends, err := services.FetchTrends(c.Context())
	if err != nil {
		fmt.Printf("Error fetching trends: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch trends: " + err.Error()})
	}
	return c.JSON(trends)
}
