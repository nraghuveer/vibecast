package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	audioMutexes sync.Map
)

func getAudioMutex(id string) *sync.Mutex {
	mu, _ := audioMutexes.LoadOrStore(id, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

func GetAudioDir(conversationID string) (string, error) {
	conversationDir, err := GetConversationDir(conversationID)
	if err != nil {
		return "", err
	}

	return filepath.Join(conversationDir, audioDirName), nil
}

func GetNextAudioIndex(conversationID string) (int, error) {
	audioDir, err := GetAudioDir(conversationID)
	if err != nil {
		return 0, err
	}

	entries, err := os.ReadDir(audioDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 1, nil
		}
		return 0, fmt.Errorf("failed to read audio directory: %w", err)
	}

	maxIndex := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if len(name) < 5 {
			continue
		}

		indexStr := name[:3]
		if name[3:] != ".wav" {
			continue
		}

		index, err := strconv.Atoi(indexStr)
		if err != nil {
			continue
		}

		if index > maxIndex {
			maxIndex = index
		}
	}

	return maxIndex + 1, nil
}

func SaveAudio(conversationID string, audioData []byte) (string, error) {
	mu := getAudioMutex(conversationID)
	mu.Lock()
	defer mu.Unlock()

	audioDir, err := GetAudioDir(conversationID)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(audioDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create audio directory: %w", err)
	}

	index, err := GetNextAudioIndex(conversationID)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%03d.wav", index)
	filepath := filepath.Join(audioDir, filename)

	if err := os.WriteFile(filepath, audioData, 0644); err != nil {
		return "", fmt.Errorf("failed to write audio file: %w", err)
	}

	return filename, nil
}

func ReadAudio(conversationID string, filename string) ([]byte, error) {
	audioDir, err := GetAudioDir(conversationID)
	if err != nil {
		return nil, err
	}

	filepath := filepath.Join(audioDir, filename)

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio file: %w", err)
	}

	return data, nil
}

func DeleteAudioFile(conversationID string, filename string) error {
	audioDir, err := GetAudioDir(conversationID)
	if err != nil {
		return err
	}

	filepath := filepath.Join(audioDir, filename)

	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete audio file: %w", err)
	}

	return nil
}

func DeleteAudioDir(conversationID string) error {
	audioDir, err := GetAudioDir(conversationID)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(audioDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete audio directory: %w", err)
	}

	return nil
}

func ListAudioFiles(conversationID string) ([]string, error) {
	audioDir, err := GetAudioDir(conversationID)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(audioDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read audio directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, entry.Name())
	}

	return files, nil
}
