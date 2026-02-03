package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/nraghuveer/vibecast/lib/config"
)

type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
	Reasoning   string        `json:"reasoning_effort,omitempty"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type chatCompletionStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

func (c *Client) chatCompletion(ctx context.Context, provider string, messages []ChatMessage) (string, error) {
	apiKey, err := config.GetProviderAPIKey(provider)
	if err != nil {
		return "", err
	}
	url, err := config.GetProviderInferenceURL(provider)
	if err != nil {
		return "", err
	}
	model, err := config.GetProviderChatModel(provider)
	if err != nil {
		return "", err
	}

	reasoningEffort := config.GetReasoningEffort()
	if provider != "openai" {
		reasoningEffort = ""
	}
	body, err := json.Marshal(chatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: 1,
		Stream:      false,
		Reasoning:   reasoningEffort,
	})
	if err != nil {
		return "", fmt.Errorf("marshal chat request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("llm chat request failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		log.Printf("llm chat error: status=%s body=%s", resp.Status, strings.TrimSpace(string(b)))
		return "", fmt.Errorf("chat completion failed: %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}

	var out chatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		log.Printf("llm chat decode error: %v", err)
		return "", fmt.Errorf("decode chat response: %w", err)
	}
	if len(out.Choices) == 0 {
		log.Printf("llm chat no choices")
		return "", errors.New("chat completion: no choices")
	}
	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}

func (c *Client) streamChatCompletion(ctx context.Context, provider string, messages []ChatMessage) (<-chan StreamEvent, error) {
	apiKey, err := config.GetProviderAPIKey(provider)
	if err != nil {
		return nil, err
	}
	url, err := config.GetProviderInferenceURL(provider)
	if err != nil {
		return nil, err
	}
	model, err := config.GetProviderChatModel(provider)
	if err != nil {
		return nil, err
	}

	reasoningEffort := config.GetReasoningEffort()
	if provider != "openai" {
		reasoningEffort = ""
	}
	body, err := json.Marshal(chatCompletionRequest{
		Model:       model,
		Messages:    messages,
		Temperature: 1,
		Stream:      true,
		Reasoning:   reasoningEffort,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal chat request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("llm chat stream request failed: %v", err)
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		log.Printf("llm chat stream error: status=%s body=%s", resp.Status, strings.TrimSpace(string(b)))
		return nil, fmt.Errorf("chat stream failed: %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}

	ch := make(chan StreamEvent, 32)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		// SSE lines are typically small, but allow some headroom.
		scanner.Buffer(make([]byte, 0, 64*1024), 256*1024)

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				log.Printf("llm chat stream canceled: %v", ctx.Err())
				ch <- StreamEvent{Done: true, Err: ctx.Err()}
				return
			default:
			}

			line := scanner.Text()
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			data = strings.TrimSpace(data)
			if data == "[DONE]" {
				ch <- StreamEvent{Done: true}
				return
			}

			var chunk chatCompletionStreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				// Ignore non-JSON lines unless context indicates fatal parsing issues.
				if !strings.Contains(data, "{") {
					continue
				}
				log.Printf("llm chat stream decode error: %v", err)
				ch <- StreamEvent{Done: true, Err: fmt.Errorf("decode stream chunk: %w", err)}
				return
			}

			for _, choice := range chunk.Choices {
				if choice.Delta.Content != "" {
					ch <- StreamEvent{Delta: choice.Delta.Content}
				}
				if choice.FinishReason != nil {
					ch <- StreamEvent{Done: true}
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("llm chat stream scanner error: %v", err)
			ch <- StreamEvent{Done: true, Err: err}
			return
		}

		// If we exit without [DONE], treat as done.
		ch <- StreamEvent{Done: true}
	}()

	return ch, nil
}
