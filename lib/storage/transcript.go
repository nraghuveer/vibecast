package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
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
