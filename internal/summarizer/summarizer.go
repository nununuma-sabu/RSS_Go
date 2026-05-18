package summarizer

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Summarizer はGemini APIとの通信を処理します
type Summarizer struct {
	client *genai.Client
}

// NewSummarizer は提供されたAPIキーを使用して新しいSummarizerインスタンスを作成します
func NewSummarizer(ctx context.Context, apiKey string) (*Summarizer, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	return &Summarizer{client: client}, nil
}

// Close は内部のクライアントを閉じます
func (s *Summarizer) Close() {
	s.client.Close()
}

// Summarize は要約のためにテキストをGemini APIに送信します
func (s *Summarizer) Summarize(ctx context.Context, text string) (string, error) {
	// 高精度な要約を行うために gemini-2.5-pro を使用します
	model := s.client.GenerativeModel("gemini-2.5-pro")
	model.SetTemperature(0.3) // より事実に基づいた要約にするためTemperatureを低く設定します

	prompt := fmt.Sprintf(`以下のRSS記事の内容（タイトルと概要、または本文）を元に、内容を3行（200文字程度）の日本語で簡潔に要約してください。
重要: 「はい、承知いたしました」などの挨拶や前置き、補足の文章は一切出力せず、要約文（3行）のみを直接出力してください。

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
