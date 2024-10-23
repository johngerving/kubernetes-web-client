package oauth

import (
	"context"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestNewConfigAndProviderFromEnv(t *testing.T) {
	normalProvider, _ := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	normalConfig := &oauth2.Config{
		RedirectURL:  "https://foo.com/callback",
		ClientID:     "oidc12345",
		ClientSecret: "secret123",
		Scopes:       []string{oidc.ScopeOpenID, "email"},
		Endpoint:     normalProvider.Endpoint(),
	}

	tests := []struct {
		description  string // Test description
		clientId     string
		clientSecret string
		callbackUrl  string
		issuer       string
		wantConfig   *oauth2.Config
		wantProvider *oidc.Provider
		wantErr      error
	}{
		{
			"Normal config",
			"oidc12345",
			"secret123",
			"https://foo.com/callback",
			"https://accounts.google.com",
			normalConfig,
			normalProvider,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Setenv("OAUTH_CLIENT_ID", test.clientId)
			t.Setenv("OAUTH_CLIENT_SECRET", test.clientSecret)
			t.Setenv("OAUTH_CALLBACK_URL", test.callbackUrl)
			t.Setenv("ISSUER", test.issuer)

			haveConfig, haveProvider, haveErr := NewConfigAndProviderFromEnv()

			if test.wantErr == nil {
				require.Nil(t, haveErr)
				require.NotNil(t, haveConfig)
				require.NotNil(t, haveProvider)

				require.Equal(t, test.wantConfig, haveConfig)
				require.Equal(t, test.wantProvider, haveProvider)
			} else {
				require.Nil(t, haveConfig)
				require.NotNil(t, haveErr)
				require.Equal(t, test.wantErr, haveErr)
			}
		})
	}
}
