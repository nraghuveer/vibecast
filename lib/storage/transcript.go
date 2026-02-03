package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/nraghuveer/vibecast/lib/models"
)

const (
	transcriptFileName = "transcript.txt"
)

var (
	transcriptMutexes sync.Map
)

func getTranscriptMutex(id string) *sync.Mutex {
	mu, _ := transcriptMutexes.LoadOrStore(id, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

func getTranscriptPath(conversationID string) (string, error) {
	conversationDir, err := GetConversationDir(conversationID)
	if err != nil {
		return "", err
	}

	return filepath.Join(conversationDir, transcriptFileName), nil
}

func CreateTranscript(conversationID string) error {
	transcriptPath, err := getTranscriptPath(conversationID)
	if err != nil {
		return err
	}

	if err := os.WriteFile(transcriptPath, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create transcript file: %w", err)
	}

	return nil
}

func AppendMessage(conversationID string, speaker string, content string) error {
	mu := getTranscriptMutex(conversationID)
	mu.Lock()
	defer mu.Unlock()

	transcriptPath, err := getTranscriptPath(conversationID)
	if err != nil {
		return err
	}

	timestamp := time.Now().Format(time.RFC3339)
	message := fmt.Sprintf("[%s] %s: %s\n", timestamp, speaker, content)

	file, err := os.OpenFile(transcriptPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open transcript file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(message); err != nil {
		return fmt.Errorf("failed to write message to transcript: %w", err)
	}

	return nil
}

func AppendMessageWithAudio(conversationID string, speaker string, content string, audioFile string) error {
	mu := getTranscriptMutex(conversationID)
	mu.Lock()
	defer mu.Unlock()

	transcriptPath, err := getTranscriptPath(conversationID)
	if err != nil {
		return err
	}

	timestamp := time.Now().Format(time.RFC3339)
	message := fmt.Sprintf("[%s] %s: %s [Audio: %s]\n", timestamp, speaker, content, audioFile)

	file, err := os.OpenFile(transcriptPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open transcript file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(message); err != nil {
		return fmt.Errorf("failed to write message to transcript: %w", err)
	}

	return nil
}

func ReadTranscript(conversationID string) (string, error) {
	transcriptPath, err := getTranscriptPath(conversationID)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(transcriptPath)
	if err != nil {
		return "", fmt.Errorf("failed to read transcript file: %w", err)
	}

	return string(data), nil
}

func TranscriptExists(conversationID string) (bool, error) {
	transcriptPath, err := getTranscriptPath(conversationID)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(transcriptPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DeleteTranscript(conversationID string) error {
	transcriptPath, err := getTranscriptPath(conversationID)
	if err != nil {
		return err
	}

	if err := os.Remove(transcriptPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete transcript file: %w", err)
	}

	return nil
}

// Message represents a parsed message from the transcript
type Message struct {
	Timestamp time.Time
	Speaker   models.SpeakerType
	Content   string
}

// LoadMessages loads all messages from the transcript file
func LoadMessages(conversationID string) ([]Message, error) {
	transcriptPath, err := getTranscriptPath(conversationID)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(transcriptPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Message{}, nil
		}
		return nil, fmt.Errorf("failed to read transcript file: %w", err)
	}

	return parseTranscript(string(data)), nil
}

// parseTranscript parses transcript content into messages
func parseTranscript(content string) []Message {
	var messages []Message
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse format: [2023-10-01T12:34:56Z] Speaker: Content
		if !strings.HasPrefix(line, "[") {
			continue
		}

		// Find closing bracket for timestamp
		closeBracket := strings.Index(line, "]")
		if closeBracket == -1 {
			continue
		}

		timestampStr := line[1:closeBracket]
		timestamp, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			continue
		}

		// Get the rest after the bracket and space
		rest := strings.TrimSpace(line[closeBracket+1:])
		if rest == "" {
			continue
		}

		// Find speaker and content
		colonIdx := strings.Index(rest, ":")
		if colonIdx == -1 {
			continue
		}

		speakerStr := strings.TrimSpace(rest[:colonIdx])
		content := strings.TrimSpace(rest[colonIdx+1:])

		// Skip audio references in content
		if idx := strings.Index(content, " [Audio:"); idx != -1 {
			content = strings.TrimSpace(content[:idx])
		}

		messages = append(messages, Message{
			Timestamp: timestamp,
			Speaker:   models.ParseSpeakerType(speakerStr),
			Content:   content,
		})
	}

	return messages
}
