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
	modelName := "gemini-2.5-flash"
	model := client.GenerativeModel(modelName)
	model.ResponseMIMEType = "application/json"

	prompt := fmt.Sprintf(`
You are an expert political analyst and social psychologist specializing in Tamil Nadu politics. 
Your task is to analyze the provided text data (Headlines, social media comments, Reddit discussions) regarding a political party.

### PHASE 1: THINKING PROCESS
Before generating the JSON, perform a deep analysis (you can output this thought process before the JSON block):
1. **Source Weighting**: Prioritize reputable news (e.g., BBC, Hindustan Times, Dinamalar) over unverified social media noise.
2. **Bias Detection**: specific political biases in the source text and neutralize them.
3. **Contextual nuance**: Differentiate between "Mockery" (trolling) and genuine "Anger". Understand TN political slang (e.g., 'Sanghi', 'Upee', 'Dravidiya Model').
4. **Aggregate Scoring**: Calculate the score based on the *weighted* evidence, not just the volume of text.

### PHASE 2: FINAL OUTPUT
Output strictly a valid JSON object.
- **Sentiment Score**: A float between -1.0 (Extreme Negative) and 1.0 (Extreme Positive).
- **Emotion**: MUST be exactly one of these: "Strong Support", "Support", "Neutral", "Disappointment", "Anger", "Hope", "Fear", "Mockery".
- **Key Topics**: Top 3-5 specific themes driving this sentiment.
- **Fact Check**: Note any identified misinformation or "None".

JSON Schema:
{
  "sentiment_score": float,
  "emotion": string,
  "key_topics": [string],
  "fact_check_notes": string
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
