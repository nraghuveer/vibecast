package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nraghuveer/vibecast/lib/config"
	"github.com/nraghuveer/vibecast/lib/data"
	"github.com/nraghuveer/vibecast/lib/db"
	"github.com/nraghuveer/vibecast/lib/logger"
)

func main() {
	// Initialize logger
	log := logger.GetInstance()
	defer log.Close()

	configPath := flag.String("config", "", "Path to config file (default: ~/.vibecast/config.yml)")
	flag.Parse()

	if _, err := config.Load(*configPath); err != nil {
		log.LogError("config_load", err)
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	log.Info("config_loaded", "path", config.GetConfigPath())

	fmt.Printf("Using config: %s\n", config.GetConfigPath())
	fmt.Printf("Database: %s\n", config.GetDBPath())
	fmt.Printf("Conversation Provider: %s\n", config.GetConversationProvider())
	fmt.Printf("Speech to Text Provider: %s\n", config.GetSpeechToTextProvider())
	fmt.Printf("Text to Speech Provider: %s\n", config.GetTextToSpeechProvider())

	log.Info("app_init", "config_path", config.GetConfigPath(), "db_path", config.GetDBPath())

	database, err := db.NewDB()
	if err != nil {
		log.LogError("database_init", err)
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()
	log.Info("database_initialized")

	data.InitializeDefaultTemplates(database)

	p := tea.NewProgram(
		NewModel(database),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.LogError("program_run", err)
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}

	log.Info("program_exited_normally")
}
