package api

import (
	"fmt"
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
	env := strings.ToLower(os.Getenv("ENV"))
	if env != "production" {
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
			return nil, fmt.Errorf("unable to load port %v: %v", portString, err)
		}
	}

	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		return nil, fmt.Errorf("API URL must be specified")
	}

	appUrl := os.Getenv("APP_URL")
	if appUrl == "" {
		return nil, fmt.Errorf("app URL must be specified")
	}

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		return nil, fmt.Errorf("domain must be specified")
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
