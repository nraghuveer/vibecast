package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/nraghuveer/vibecast/lib/config"
)

type ttsRequest struct {
	Model          string `json:"model"`
	Voice          string `json:"voice"`
	Input          string `json:"input"`
	ResponseFormat string `json:"response_format,omitempty"`
}

func (c *Client) openAITTS(ctx context.Context, provider, voice, text string) ([]byte, error) {
	apiKey, err := config.GetProviderAPIKey(provider)
	if err != nil {
		return nil, err
	}
	url, err := config.GetProviderTTSURL(provider)
	if err != nil {
		return nil, err
	}
	model, err := config.GetProviderTTSModel(provider)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(url) == "" {
		return nil, fmt.Errorf("tts url not configured for provider %s", provider)
	}
	if strings.TrimSpace(model) == "" {
		return nil, fmt.Errorf("tts model not configured for provider %s", provider)
	}

	body, err := json.Marshal(ttsRequest{
		Model:          model,
		Voice:          voice,
		Input:          text,
		ResponseFormat: "wav",
	})
	if err != nil {
		return nil, fmt.Errorf("marshal tts request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("tts request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		log.Printf("tts error: status=%s body=%s", resp.Status, strings.TrimSpace(string(b)))
		return nil, fmt.Errorf("tts failed: %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}

	audio, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("tts read error: %v", err)
		return nil, err
	}
	if len(audio) == 0 {
		log.Printf("tts returned empty audio")
		return nil, fmt.Errorf("tts returned empty audio")
	}
	return audio, nil
}
