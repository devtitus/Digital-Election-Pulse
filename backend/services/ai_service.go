package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIAnalysisResult struct {
	SentimentScore float64  `json:"sentiment_score"`
	Emotion        string   `json:"emotion"`
	KeyTopics      []string `json:"key_topics"`
	FactCheckNotes string   `json:"fact_check_notes"`
}

func AnalyzeSentiment(ctx context.Context, textData string) (*AIAnalysisResult, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	defer client.Close()

	// Try gemini-1.5-flash-001 first, then fallback to gemini-pro if needed
	modelName := "gemini-robotics-er-1.5-preview"
	model := client.GenerativeModel(modelName)
	model.ResponseMIMEType = "application/json"

	prompt := fmt.Sprintf(`
You are an expert political analyst and social psychologist specializing in Tamil Nadu politics. 
Your task is to analyze the provided text data (Headlines, social media comments, Reddit discussions) regarding a political party.

Guidelines:
1. **Bias Mitigation**: Remain strictly objective. Audit your analysis for any partisan bias.
2. **Fact-Checking**: Cross-reference claims in the text. If a claim is demonstrably false or misinformation, lower the sentiment score and note it.
3. **Emotional Quotient (EQ)**: Go beyond basic sentiment. Identify deep emotional drivers (e.g., "Resentment regarding flood relief" is deeper than just "Anger").
4. **Context Awareness**: Understand Tamil internet slang and cultural nuances (e.g., 'Sanghi', 'Upee', '200 rs', 'Vadinokkan', 'Thalaiva', 'Mass').
5. **Spam Filtering**: Ignore repetitive bot-like comments or irrelevant noise.

Output strictly a valid JSON object with this schema. Do NOT use markdown code blocks or any other formatting. Only the raw JSON.
{
  "sentiment_score": float (-1.0 to 1.0),
  "emotion": string (One of: Anger, Hope, Mockery, Support, Fear, Indifference, Pride),
  "key_topics": array of strings (Top 3-5 recurring themes, specific and specific),
  "fact_check_notes": string (Brief note on any major misinformation detected, or "None")
}

Data to Analyze:
%s
`, textData)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("gemini generation failed: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	// Extract text
	part := resp.Candidates[0].Content.Parts[0]
	txt, ok := part.(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	// Extract JSON from response (robust method)
	rawText := strings.TrimSpace(string(txt))

	start := strings.Index(rawText, "{")
	end := strings.LastIndex(rawText, "}")

	if start == -1 || end == -1 || start > end {
		return nil, fmt.Errorf("invalid response format: could not find JSON object. Raw: %s", rawText)
	}

	jsonStr := rawText[start : end+1]

	var result AIAnalysisResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON response: %w. Raw: %s", err, jsonStr)
	}

	return &result, nil
}
