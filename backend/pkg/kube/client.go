package kube

import (
	"encoding/base64"
	"log"

	"github.com/johngerving/kubernetes-web-client/backend/pkg/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubeClient struct {
	clientset *kubernetes.Clientset
	Namespace string
}

// NewKubeClient creates a KubeClient using a config.KubeConfig
func NewKubeClient(cfg config.KubeConfig) (*KubeClient, error) {
	// Get CA data from environment variable and base64 decode it
	caData, err := base64.StdEncoding.DecodeString(cfg.Cert)
	if err != nil {
		log.Fatalf("Error decoding cert data: %v\n", caData)
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
		log.Fatalf("Error creating Kubernetes clientset: %v", err)
	}

	kubeClient := KubeClient{
		clientset: clientset,
		Namespace: cfg.Namespace,
	}

	return &kubeClient, nil
}
