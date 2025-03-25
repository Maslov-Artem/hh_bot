package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	JobAPIURL    string
	JobAPIKey    string
	LLMAPIURL    string
	LLMAPIKey    string
	DatabaseURL  string
	Model        string
	SystemPrompt string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file\n", err)
	}
	return &Config{
		JobAPIURL:    os.Getenv("JOB_API_URL"),
		JobAPIKey:    os.Getenv("JOB_API_KEY"),
		LLMAPIURL:    os.Getenv("LLM_API_URL"),
		LLMAPIKey:    os.Getenv("LLM_API_KEY"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		Model:        os.Getenv("MODEL"),
		SystemPrompt: os.Getenv("SYSTEM_PROMPT"),
	}
}
