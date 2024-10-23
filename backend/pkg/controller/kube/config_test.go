package kube

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigFromEnv(t *testing.T) {
	tests := []struct {
		description string // Test description
		host        string
		port        string
		token       string
		cert        string
		namespace   string
		wantConfig  *KubeConfig
		wantErr     error
	}{
		{"Normal config", "127.0.0.1", "6443", "abcdefg", "hijklmnop", "default", &KubeConfig{"127.0.0.1", "6443", "abcdefg", "hijklmnop", "default"}, nil},
		{"Missing KUBERNETES_SERVICE_HOST variable", "", "6443", "abcdefg", "hijklmnop", "default", nil, fmt.Errorf("could not retrieve Kubernetes service host")},
		{"Missing KUBERNETES_SERVICE_PORT_HTTPS variable", "127.0.0.1", "", "abcdefg", "hijklmnop", "default", nil, fmt.Errorf("could not retrieve Kubernetes service port")},
		{"Missing KUBE_TOKEN variable", "127.0.0.1", "6443", "", "hijklmnop", "default", nil, fmt.Errorf("kubernetes token must be specified")},
		{"Missing KUBE_CERT variable", "127.0.0.1", "6443", "abcdefg", "", "default", nil, fmt.Errorf("kubernetes CA cert must be specified")},
		{"Missing POD_NAMESPACE variable", "127.0.0.1", "6443", "abcdefg", "hijklmnop", "", nil, fmt.Errorf("could not retrieve Kubernetes namespace")},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Setenv("KUBERNETES_SERVICE_HOST", test.host)
			t.Setenv("KUBERNETES_SERVICE_PORT_HTTPS", test.port)
			t.Setenv("KUBE_TOKEN", test.token)
			t.Setenv("KUBE_CERT", test.cert)
			t.Setenv("POD_NAMESPACE", test.namespace)

			haveConfig, haveErr := NewKubeConfigFromEnv()

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
