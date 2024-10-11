package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type config struct {
	Env         string
	Port        int
	ApiUrl      string
	AppUrl      string
	Domain      string
	OAuthUrl    string
	OAuthConfig oauth2.Config
}

var AppConfig config

func NewConfigFromEnv() error {
	// Load environment variables from '.env' file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("ENV")
	if strings.ToLower(env) != "production" {
		env = "development"
	}

	portString := os.Getenv("PORT")
	var port int
	if portString == "" {
		port = 8080
	} else {
		port, err = strconv.Atoi(portString)

		if err != nil {
			log.Fatalf("Error loading port %v: %v", portString, err)
		}
	}

	oauthClientId := os.Getenv("OAUTH_CLIENT_ID")
	if oauthClientId == "" {
		log.Fatalf("Error: OAuth client ID must be specified")
	}

	oauthClientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
	if oauthClientSecret == "" {
		log.Fatalf("Error: OAuth client secret must be specified")
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

	oauthUrl := os.Getenv("OAUTH_URL")
	if oauthUrl == "" {
		log.Fatalf("Error: OAuth URL must be specified")
	}

	oauthCallbackUrl := os.Getenv("OAUTH_CALLBACK_URL")
	if oauthCallbackUrl == "" {
		log.Fatalf("Error: OAuth callback URL must be specified")
	}

	// Create OAuth config from provided variables
	oauthConfig := oauth2.Config{
		RedirectURL:  oauthCallbackUrl,
		ClientID:     oauthClientId,
		ClientSecret: oauthClientSecret,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	// Create AppConfig, including the oauthConfig
	AppConfig = config{
		Env:         env,
		Port:        port,
		ApiUrl:      apiUrl,
		AppUrl:      appUrl,
		Domain:      domain,
		OAuthUrl:    oauthUrl,
		OAuthConfig: oauthConfig,
	}

	return nil
}
