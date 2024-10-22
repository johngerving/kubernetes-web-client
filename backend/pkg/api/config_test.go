package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigFromEnv(t *testing.T) {
	tests := []struct {
		description string // Test description
		env         string
		port        string
		apiUrl      string
		appUrl      string
		domain      string
		wantConfig  *Config
		wantErr     error
	}{
		{"Normal config", "Production", "8090", "foo.com/api", "foo.com", "foo.com", &Config{"production", 8090, "foo.com/api", "foo.com", "foo.com"}, nil},
		{"Missing ENV variable", "", "8090", "foo.com/api", "foo.com", "foo.com", &Config{"development", 8090, "foo.com/api", "foo.com", "foo.com"}, nil},
		{"Missing PORT variable", "production", "", "foo.com/api", "foo.com", "foo.com", &Config{"production", 8080, "foo.com/api", "foo.com", "foo.com"}, nil},
		{"Missing API_URL variable", "production", "8090", "", "foo.com", "foo.com", nil, fmt.Errorf("API URL must be specified")},
		{"Missing APP_URL variable", "production", "8090", "foo.com/api", "", "foo.com", nil, fmt.Errorf("app URL must be specified")},
		{"Missing DOMAIN variable", "production", "8090", "foo.com/api", "foo.com", "", nil, fmt.Errorf("domain must be specified")},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Setenv("ENV", test.env)
			t.Setenv("PORT", test.port)
			t.Setenv("API_URL", test.apiUrl)
			t.Setenv("APP_URL", test.appUrl)
			t.Setenv("DOMAIN", test.domain)

			haveConfig, haveErr := NewConfigFromEnv()

			if test.wantErr == nil {
				require.Nil(t, haveErr)
				require.NotNil(t, haveConfig)
				require.Equal(t, test.wantConfig, haveConfig)
			} else {
				require.Nil(t, haveConfig)
				require.NotNil(t, haveErr)
				require.Equal(t, test.wantErr, haveErr)
			}
		})
	}
}
