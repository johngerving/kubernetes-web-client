package kube

import (
	"encoding/base64"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubeController struct {
	clientset *kubernetes.Clientset
	Namespace string
}

// NewKubeController creates a KubeControl using a kube.KubeConfig
func NewKubeController(cfg *KubeConfig) (*KubeController, error) {
	// Get CA data from environment variable and base64 decode it
	caData, err := base64.StdEncoding.DecodeString(cfg.Cert)
	if err != nil {
		return nil, fmt.Errorf("unable to decode cert data: %v", caData)
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

	kubeClient := &KubeController{
		clientset: clientset,
		Namespace: cfg.Namespace,
	}

	return kubeClient, nil
}

func (c KubeController) GetWorkspacePodStatus(username string)    {}
func (c KubeController) GetWorkspaceVolumeStatus(username string) {}
func (c KubeController) CreateWorkspacePod(username string)       {}
func (c KubeController) CreateWorkspaceVolume(username string)    {}
