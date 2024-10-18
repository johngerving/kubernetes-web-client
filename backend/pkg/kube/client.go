package kube

import (
	"encoding/base64"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	clientset *kubernetes.Clientset
	Namespace string
}

// NewClient creates a Client using a kube.Config
func NewClient(cfg *Config) (*Client, error) {
	// Get CA data from environment variable and base64 decode it
	caData, err := base64.StdEncoding.DecodeString(cfg.Cert)
	if err != nil {
		return nil, fmt.Errorf("unable to decode cert data: %v\n", caData)
	}

	// Create a rest config
	config := &rest.Config{
		Host:        "https://" + cfg.Host + ":" + cfg.Port,
		BearerToken: cfg.Token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: false,
			CAData:   caData, // Pass in CA data for TLS verification
		},
	}

	// Create a clientset from the config
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create Kubernetes clientset: %v", err)
	}

	kubeClient := &Client{
		clientset: clientset,
		Namespace: cfg.Namespace,
	}

	return kubeClient, nil
}
