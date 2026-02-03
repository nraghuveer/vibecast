package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nraghuveer/vibecast/lib/config"
	"github.com/nraghuveer/vibecast/lib/data"
	"github.com/nraghuveer/vibecast/lib/db"
)

func main() {
	configPath := flag.String("config", "", "Path to config file (default: ~/.vibecast/config.yml)")
	flag.Parse()

	if _, err := config.Load(*configPath); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Using config: %s\n", config.GetConfigPath())
	fmt.Printf("Database: %s\n", config.GetDBPath())
	fmt.Printf("Conversation Provider: %s\n", config.GetConversationProvider())
	fmt.Printf("Speech to Text Provider: %s\n", config.GetSpeechToTextProvider())
	fmt.Printf("Text to Speech Provider: %s\n", config.GetTextToSpeechProvider())

	if err := db.Init(); err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	data.InitializeDefaultTemplates()

	p := tea.NewProgram(
		NewModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
