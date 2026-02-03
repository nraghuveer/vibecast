package llm

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	prompts    *PromptLoader
}

func New() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 90 * time.Second},
		prompts:    NewPromptLoader("prompts"),
	}
}

// StreamGuestResponse streams the next guest reply for the given conversation.
// It emits incremental deltas suitable for real-time UI updates.
func (c *Client) StreamGuestResponse(ctx context.Context, provider, persona, topic string, history []ChatMessage) (<-chan StreamEvent, error) {
	sys, err := c.prompts.RenderFile("system_prompt.txt", struct {
		Persona string
		Topic   string
	}{
		Persona: persona,
		Topic:   topic,
	})
	if err != nil {
		return nil, err
	}

	msgs := make([]ChatMessage, 0, len(history)+1)
	msgs = append(msgs, ChatMessage{Role: "system", Content: sys})
	msgs = append(msgs, history...)

	return c.streamChatCompletion(ctx, provider, msgs)
}

// prepareTextForSpeech rewrites text into natural, speakable dialogue.
// This is useful before feeding text into a TTS engine.
func (c *Client) prepareTextForSpeech(ctx context.Context, provider, persona, topic, voice, text string) (string, error) {
	prompt, err := c.prompts.RenderFile("text_to_speech.txt", struct {
		Persona string
		Topic   string
		Voice   string
		Text    string
	}{
		Persona: persona,
		Topic:   topic,
		Voice:   voice,
		Text:    text,
	})
	if err != nil {
		return "", err
	}

	// Keep this as a single system message. The template includes the input.
	msgs := []ChatMessage{{Role: "system", Content: prompt}}
	out, err := c.chatCompletion(ctx, provider, msgs)
	if err != nil {
		return "", err
	}
	if out == "" {
		return "", fmt.Errorf("empty tts prep output")
	}
	return out, nil
}

// SynthesizeGuestSpeech converts guest text into audio.
// It first normalizes text into natural, human-like speech, then calls the TTS endpoint.
// Returns the synthesized audio bytes and the speakable text actually sent to TTS.
func (c *Client) SynthesizeGuestSpeech(ctx context.Context, prepProvider, ttsProvider, persona, topic, voice, text string) ([]byte, string, error) {
	speakable := text
	if strings.TrimSpace(prepProvider) != "" {
		if prepared, err := c.prepareTextForSpeech(ctx, prepProvider, persona, topic, voice, text); err == nil && strings.TrimSpace(prepared) != "" {
			speakable = prepared
		}
	}

	audio, err := c.openAITTS(ctx, ttsProvider, voice, speakable)
	if err != nil {
		return nil, speakable, err
	}
	return audio, speakable, nil
}
