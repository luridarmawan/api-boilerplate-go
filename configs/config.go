package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIName    string
	APIDescription string
	APIVersion string
	BaseURL    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	ServerPort string
	DocFilter  string

	// AI Configuration
	AIBaseURL string
	AIAPIKey  string
	AITimeout string

	// Build info
	Version   string
	GitCommit string
	BuildDate string
}

// These variables are injected at build time using -ldflags
var (
	Version   string
	GitCommit string
	BuildDate string
)

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		APIName:     getEnv("API_NAME", "My API Name"),
		APIDescription:     getEnv("API_DESCRIPTION", "My API Description"),
		APIVersion:     getEnv("API_VERSION", "0.0.0"),
		BaseURL:     getEnv("BASEURL", "localhost:" + getEnv("SERVER_PORT", "3000")),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "my_api_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		ServerPort: getEnv("SERVER_PORT", "3000"),
		DocFilter:  getEnv("API_DOC_FILTER", ""),

		// AI Configuration
		AIBaseURL: getEnv("AI_BASE_URL", "https://api.openai.com/v1"),
		AIAPIKey:  getEnv("AI_API_KEY", ""),
		AITimeout: getEnv("AI_TIMEOUT", "30"),

		// Inject build-time values
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}