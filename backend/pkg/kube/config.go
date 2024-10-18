package kube

import (
	"fmt"
	"os"
)

type Config struct {
	Host      string
	Port      string
	Token     string
	Cert      string
	Namespace string
}

func NewConfigFromEnv() (*Config, error) {
	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	if host == "" {
		return nil, fmt.Errorf("could not retrieve Kubernetes service host")
	}

	port := os.Getenv("KUBERNETES_SERVICE_PORT_HTTPS")
	if port == "" {
		return nil, fmt.Errorf("could not retrieve Kubernetes service port")
	}

	token := os.Getenv("KUBE_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("Kubernetes token must be specified")
	}

	cert := os.Getenv("KUBE_CERT")
	if cert == "" {
		return nil, fmt.Errorf("Kubernetes CA cert must be specified")
	}

	namespace := os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		return nil, fmt.Errorf("could not retrieve Kubernetes namespace")
	}

	cfg := &Config{
		Host:      host,
		Port:      port,
		Token:     token,
		Cert:      cert,
		Namespace: namespace,
	}

	return cfg, nil
}
