package api

import (
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Environment string
	Port        int
	BackendURL  string
	FrontendURL string
	Domain      string
}

// NewConfigFromEnv reads in environment variables and returns
// a Config struct instance. If an environment variable is not found,
// an error occurs.
func NewConfigFromEnv() (*Config, error) {
	// Load environment variables
	env := os.Getenv("ENV")
	if strings.ToLower(env) != "production" {
		env = "development"
	}

	portString := os.Getenv("PORT")
	var port int
	var err error
	if portString == "" {
		port = 8080
	} else {
		port, err = strconv.Atoi(portString)

		if err != nil {
			log.Fatalf("Error loading port %v: %v", portString, err)
		}
	}

	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		log.Fatalf("Error: API URL must be specified")
	}

	appUrl := os.Getenv("APP_URL")
	if appUrl == "" {
		log.Fatalf("Error: App URL must be specified")
	}

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		log.Fatalf("Error: Domain must be specified")
	}

	// Create config, including the oauthConfig
	cfg := Config{
		Environment: env,
		Port:        port,
		BackendURL:  apiUrl,
		FrontendURL: appUrl,
		Domain:      domain,
	}

	return &cfg, nil
}
