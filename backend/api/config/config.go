package config

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

type Config struct {
	Env         string
	Port        int
	ApiUrl      string
	AppUrl      string
	Domain      string
	OAuthUrl    string
	OAuthConfig oauth2.Config
	Provider    *oidc.Provider
}

// NewConfigFromEnv reads in environment variables and returns
// a Config struct instance. If an environment variable is not found,
// an error occurs.
func NewConfigFromEnv() (Config, error) {
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

	oauthCallbackUrl := os.Getenv("OAUTH_CALLBACK_URL")
	if oauthCallbackUrl == "" {
		log.Fatalf("Error: OAuth callback URL must be specified")
	}

	issuer := os.Getenv("ISSUER")
	if issuer == "" {
		log.Fatalf("Error: Issuer URL must be specified")
	}

	// Create OIDC provider
	provider, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		log.Fatalf("Error creating OIDC provider")
	}

	// Create OpenID Connect aware OAuth config from provided variables
	oauthConfig := oauth2.Config{
		RedirectURL:  oauthCallbackUrl,
		ClientID:     oauthClientId,
		ClientSecret: oauthClientSecret,
		Scopes:       []string{oidc.ScopeOpenID, "email"}, // "openid" is a required scope for OpenID Connect flows
		Endpoint:     provider.Endpoint(),
	}

	// Create config, including the oauthConfig
	cfg := Config{
		Env:         env,
		Port:        port,
		ApiUrl:      apiUrl,
		AppUrl:      appUrl,
		Domain:      domain,
		OAuthConfig: oauthConfig,
		Provider:    provider,
	}

	return cfg, nil
}
