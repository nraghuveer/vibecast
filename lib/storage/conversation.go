package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	conversationsDirName = "conversations"
	audioDirName         = "audio"
)

func GetVibecastDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".vibecast"), nil
}

func GetConversationsDir() (string, error) {
	vibeDir, err := GetVibecastDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(vibeDir, conversationsDirName), nil
}

func CreateConversationDir(id string) (string, error) {
	conversationsDir, err := GetConversationsDir()
	if err != nil {
		return "", err
	}

	conversationDir := filepath.Join(conversationsDir, id)

	if err := os.MkdirAll(conversationDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create conversation directory: %w", err)
	}

	audioDir := filepath.Join(conversationDir, audioDirName)
	if err := os.MkdirAll(audioDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create audio directory: %w", err)
	}

	return conversationDir, nil
}

func GetConversationDir(id string) (string, error) {
	conversationsDir, err := GetConversationsDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(conversationsDir, id), nil
}

func ConversationExists(id string) (bool, error) {
	conversationDir, err := GetConversationDir(id)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(conversationDir)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DeleteConversationDir(id string) error {
	conversationDir, err := GetConversationDir(id)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(conversationDir); err != nil {
		return fmt.Errorf("failed to delete conversation directory: %w", err)
	}

	return nil
}

func DeleteAllConversationsDirs() error {
	conversationsDir, err := GetConversationsDir()
	if err != nil {
		return err
	}

	if err := os.RemoveAll(conversationsDir); err != nil {
		return fmt.Errorf("failed to delete conversations directory: %w", err)
	}

	return nil
}
