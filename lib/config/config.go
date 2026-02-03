package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	General   GeneralConfig             `yaml:"general"`
	AI        AIConfig                  `yaml:"ai"`
	UI        UIConfig                  `yaml:"ui"`
	Providers map[string]ProviderConfig `yaml:"providers"`
}

type GeneralConfig struct {
	DBPath string `yaml:"db_path"`
}

type UIConfig struct {
	ShowTranscripts bool           `yaml:"show_transcripts"`
	TranscriptSide  TranscriptSide `yaml:"transcript_side"`
	TranscriptWidth int            `yaml:"transcript_width"`
	Wave            WaveConfig     `yaml:"wave"`
}

type WaveConfig struct {
	Phase     float64 `yaml:"phase"`
	Frequency float64 `yaml:"frequency"`
	Amplitude int     `yaml:"amplitude"`
}

type TranscriptSide string

const (
	TranscriptSideLeft  TranscriptSide = "left"
	TranscriptSideRight TranscriptSide = "right"
)

type AIConfig struct {
	ConversationProvider string `yaml:"conversation_provider"`
	SpeechToTextProvider string `yaml:"speech_to_text"`
	TextToSpeechProvider string `yaml:"text_to_speech"`
	ReasoningEffort      string `yaml:"reasoning_effort"`
}

type ProviderConfig struct {
	ChatModel    string `yaml:"chat_model"`
	STTModel     string `yaml:"stt_model"`
	TTSModel     string `yaml:"tts_model"`
	InferenceURL string `yaml:"inference_url"`
	STTURL       string `yaml:"stt_url"`
	TTSURL       string `yaml:"tts_url"`
	APIKey       string `yaml:"api_key"`
	IsEnvVar     bool   `yaml:"is_env_var"`
}

const (
	defaultConfigDir  = ".vibecast"
	defaultConfigFile = "config.yml"
	defaultDBFile     = "data.sqlite"
	defaultProvider   = "groq"
)

var (
	globalConfig *Config
	configPath   string
)

func Load(configFilePath string) (*Config, error) {
	if configFilePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		configFilePath = filepath.Join(homeDir, defaultConfigDir, defaultConfigFile)
	}

	configPath = configFilePath

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := createDefaultConfig()
			if err := Save(cfg, configFilePath); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
			globalConfig = &cfg
			return globalConfig, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	cfg.setDefaults()
	globalConfig = &cfg
	return globalConfig, nil
}

func createDefaultConfig() Config {
	homeDir, _ := os.UserHomeDir()
	defaultDBPath := filepath.Join(homeDir, defaultConfigDir, defaultDBFile)

	return Config{
		General: GeneralConfig{
			DBPath: defaultDBPath,
		},
		AI: AIConfig{
			ConversationProvider: defaultProvider,
			SpeechToTextProvider: defaultProvider,
			TextToSpeechProvider: defaultProvider,
			ReasoningEffort:      "",
		},
		UI: UIConfig{
			ShowTranscripts: true,
			TranscriptSide:  TranscriptSideRight,
			TranscriptWidth: 40,
			Wave: WaveConfig{
				Phase:     0,
				Frequency: 0.05,
				Amplitude: 3,
			},
		},
		Providers: map[string]ProviderConfig{
			"groq": {
				ChatModel:    "llama-3.3-70b-versatile",
				STTModel:     "",
				TTSModel:     "",
				InferenceURL: "https://api.groq.com/openai/v1/chat/completions",
				STTURL:       "",
				TTSURL:       "",
				APIKey:       "",
				IsEnvVar:     true,
			},
			"openai": {
				ChatModel:    "gpt-4o",
				STTModel:     "whisper-1",
				TTSModel:     "tts-1",
				InferenceURL: "https://api.openai.com/v1/chat/completions",
				STTURL:       "https://api.openai.com/v1/audio/transcriptions",
				TTSURL:       "https://api.openai.com/v1/audio/speech",
				APIKey:       "",
				IsEnvVar:     true,
			},
		},
	}
}

func (c *Config) setDefaults() {
	if c.General.DBPath == "" {
		homeDir, _ := os.UserHomeDir()
		c.General.DBPath = filepath.Join(homeDir, defaultConfigDir, defaultDBFile)
	}

	if c.AI.ConversationProvider == "" {
		c.AI.ConversationProvider = defaultProvider
	}

	if c.AI.SpeechToTextProvider == "" {
		c.AI.SpeechToTextProvider = defaultProvider
	}

	if c.AI.TextToSpeechProvider == "" {
		c.AI.TextToSpeechProvider = defaultProvider
	}

	if c.UI.TranscriptSide == "" {
		c.UI.TranscriptSide = TranscriptSideRight
	}

	if c.UI.TranscriptWidth == 0 {
		c.UI.TranscriptWidth = 40
	}

	c.UI.Wave.Phase = 0
	if c.UI.Wave.Frequency == 0 {
		c.UI.Wave.Frequency = 0.05
	}
	if c.UI.Wave.Amplitude == 0 {
		c.UI.Wave.Amplitude = 3
	}

	if c.Providers == nil {
		c.Providers = make(map[string]ProviderConfig)
	}

	if _, exists := c.Providers["groq"]; !exists {
		c.Providers["groq"] = ProviderConfig{
			ChatModel:    "llama-3.3-70b-versatile",
			STTModel:     "",
			TTSModel:     "",
			InferenceURL: "https://api.groq.com/openai/v1/chat/completions",
			STTURL:       "",
			TTSURL:       "",
			APIKey:       "",
			IsEnvVar:     true,
		}
	}

	if _, exists := c.Providers["openai"]; !exists {
		c.Providers["openai"] = ProviderConfig{
			ChatModel:    "gpt-4o",
			STTModel:     "whisper-1",
			TTSModel:     "tts-1",
			InferenceURL: "https://api.openai.com/v1/chat/completions",
			STTURL:       "https://api.openai.com/v1/audio/transcriptions",
			TTSURL:       "https://api.openai.com/v1/audio/speech",
			APIKey:       "",
			IsEnvVar:     true,
		}
	}
}

func Save(cfg Config, configFilePath string) error {
	configDir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	tempPath := configFilePath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary config file: %w", err)
	}

	if err := os.Rename(tempPath, configFilePath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to replace config file: %w", err)
	}

	if err := removeSwapFiles(configDir); err != nil {
		fmt.Printf("Warning: failed to remove swap files: %v\n", err)
	}

	return nil
}

func removeSwapFiles(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".swp") ||
			strings.HasSuffix(name, ".swo") ||
			strings.HasSuffix(name, "~") {
			path := filepath.Join(dir, name)
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove swap file %s: %w", path, err)
			}
		}
	}

	return nil
}

func Get() *Config {
	return globalConfig
}

func GetDBPath() string {
	if globalConfig != nil && globalConfig.General.DBPath != "" {
		return globalConfig.General.DBPath
	}
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, defaultConfigDir, defaultDBFile)
}

func GetConversationProvider() string {
	if globalConfig != nil && globalConfig.AI.ConversationProvider != "" {
		return globalConfig.AI.ConversationProvider
	}
	return defaultProvider
}

func GetSpeechToTextProvider() string {
	if globalConfig != nil && globalConfig.AI.SpeechToTextProvider != "" {
		return globalConfig.AI.SpeechToTextProvider
	}
	return defaultProvider
}

func GetTextToSpeechProvider() string {
	if globalConfig != nil && globalConfig.AI.TextToSpeechProvider != "" {
		return globalConfig.AI.TextToSpeechProvider
	}
	return defaultProvider
}

func GetReasoningEffort() string {
	if globalConfig != nil {
		return strings.TrimSpace(globalConfig.AI.ReasoningEffort)
	}
	return ""
}

func GetProviderConfig(provider string) (*ProviderConfig, error) {
	if globalConfig == nil {
		return nil, fmt.Errorf("config not loaded")
	}

	cfg, exists := globalConfig.Providers[provider]
	if !exists {
		return nil, fmt.Errorf("provider %s not found in config", provider)
	}

	return &cfg, nil
}

func GetConfigPath() string {
	return configPath
}

func GetProviderAPIKey(provider string) (string, error) {
	cfg, err := GetProviderConfig(provider)
	if err != nil {
		return "", err
	}

	if cfg.IsEnvVar {
		envVar := strings.ToUpper(provider + "_api_key")
		apiKey := os.Getenv(envVar)
		if apiKey == "" {
			return "", fmt.Errorf("environment variable %s is not set", envVar)
		}
		return apiKey, nil
	}

	if cfg.APIKey == "" {
		envVar := strings.ToUpper(provider + "_api_key")
		apiKey := os.Getenv(envVar)
		if apiKey != "" {
			return apiKey, nil
		}
		return "", fmt.Errorf("api_key not configured and environment variable %s not set", envVar)
	}

	return cfg.APIKey, nil
}

func GetProviderModel(provider string) (string, error) {
	cfg, err := GetProviderConfig(provider)
	if err != nil {
		return "", err
	}
	return cfg.ChatModel, nil
}

func GetProviderChatModel(provider string) (string, error) {
	cfg, err := GetProviderConfig(provider)
	if err != nil {
		return "", err
	}
	return cfg.ChatModel, nil
}

func GetProviderSTTModel(provider string) (string, error) {
	cfg, err := GetProviderConfig(provider)
	if err != nil {
		return "", err
	}
	return cfg.STTModel, nil
}

func GetProviderTTSModel(provider string) (string, error) {
	cfg, err := GetProviderConfig(provider)
	if err != nil {
		return "", err
	}
	return cfg.TTSModel, nil
}

func GetProviderSTTURL(provider string) (string, error) {
	cfg, err := GetProviderConfig(provider)
	if err != nil {
		return "", err
	}
	return cfg.STTURL, nil
}

func GetProviderTTSURL(provider string) (string, error) {
	cfg, err := GetProviderConfig(provider)
	if err != nil {
		return "", err
	}
	return cfg.TTSURL, nil
}

func GetProviderInferenceURL(provider string) (string, error) {
	cfg, err := GetProviderConfig(provider)
	if err != nil {
		return "", err
	}
	return cfg.InferenceURL, nil
}

func GetUIConfig() UIConfig {
	if globalConfig != nil {
		return globalConfig.UI
	}
	return UIConfig{
		ShowTranscripts: true,
		TranscriptSide:  TranscriptSideRight,
		TranscriptWidth: 40,
		Wave: WaveConfig{
			Phase:     0,
			Frequency: 0.05,
			Amplitude: 3,
		},
	}
}
