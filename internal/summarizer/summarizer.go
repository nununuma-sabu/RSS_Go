package summarizer

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Summarizer handles communication with Gemini API
type Summarizer struct {
	client *genai.Client
}

// NewSummarizer creates a new Summarizer instance using the provided API key
func NewSummarizer(ctx context.Context, apiKey string) (*Summarizer, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	return &Summarizer{client: client}, nil
}

// Close closes the underlying client
func (s *Summarizer) Close() {
	s.client.Close()
}

// Summarize sends the text to gemini-1.5-pro for summarization
func (s *Summarizer) Summarize(ctx context.Context, text string) (string, error) {
	// Use gemini-2.5-pro for high accuracy summarization
	model := s.client.GenerativeModel("gemini-2.5-pro")
	model.SetTemperature(0.3) // Lower temperature for more factual summarization

	prompt := fmt.Sprintf(`以下のRSS記事の内容（タイトルと概要、または本文）を元に、内容を3行（200文字程度）の日本語で簡潔に要約してください。

---
%s
---`, text)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no summary generated")
	}

	part := resp.Candidates[0].Content.Parts[0]
	if textPart, ok := part.(genai.Text); ok {
		return strings.TrimSpace(string(textPart)), nil
	}

	return "", fmt.Errorf("unexpected response format")
}
