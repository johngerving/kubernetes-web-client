package controller

import (
	"fmt"
	"os"
	"strings"

	"github.com/johngerving/kubernetes-web-client/backend/pkg/controller/kube"
	_ "github.com/joho/godotenv/autoload"
)

type Controller interface {
	GetWorkspacePodStatus(username string)
	GetWorkspaceVolumeStatus(username string)
	CreateWorkspacePod(username string)
	CreateWorkspaceVolume(username string)
}

// NewControllerFromEnv creates a new Controller interface instance
// from environment variables.
func NewControllerFromEnv() (Controller, error) {
	// Get environment variables
	clusterType := strings.ToLower(os.Getenv("CLUSTER_TYPE"))

	if clusterType == "" {
		return nil, fmt.Errorf("cluster type must be specified")
	}

	// Check which type of cluster is being used
	if clusterType == "kubernetes" {
		// If Kubernetes cluster, create a new Kubernetes Controller
		cfg, err := kube.NewKubeConfigFromEnv()
		if err != nil {
			return nil, err
		}

		controller, err := kube.NewKubeController(cfg)
		if err != nil {
			return nil, err
		}

		return controller, nil
	}

	return nil, fmt.Errorf("invalid cluster type")
}
