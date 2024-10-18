package oauth

import (
	"context"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

func NewConfigAndProviderFromEnv() (*oauth2.Config, *oidc.Provider, error) {
	// Get environment variables

	clientId := os.Getenv("OAUTH_CLIENT_ID")
	if clientId == "" {
		return nil, nil, fmt.Errorf("OAuth client secret must be specified")
	}

	clientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, nil, fmt.Errorf("OAuth client secret must be specified")
	}

	callback := os.Getenv("OAUTH_CALLBACK_URL")
	if callback == "" {
		return nil, nil, fmt.Errorf("OAuth callback URL must be specified")
	}

	issuer := os.Getenv("ISSUER")
	if issuer == "" {
		return nil, nil, fmt.Errorf("issuer URL must be specified")
	}

	// Create OIDC provider
	provider, err := oidc.NewProvider(context.Background(), issuer)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create OIDC provider: %v", err)
	}

	// Create OpenID Connect aware OAuth config from provided variables
	config := &oauth2.Config{
		RedirectURL:  callback,
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{oidc.ScopeOpenID, "email"}, // "openid" is a required scope for OpenID Connect flows
		Endpoint:     provider.Endpoint(),
	}

	return config, provider, nil
}
