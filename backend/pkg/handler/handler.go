package handler

import (
	"github.com/alexedwards/scs/v2"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/config"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/database/repository"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/kube"
)

// Struct to hold app information - config, stores, etc.
type Handler struct {
	AppConfig      *config.Config
	SessionManager *scs.SessionManager
	Repository     *repository.Queries
	KubeClient     *kube.KubeClient
}

// NewHandler takes in a config.Config struct instance, a scs.SessionManager pointer, and a
// repository.Queries pointer and returns a pointer to a Handler struct instance.
func NewHandler(appConfig *config.Config, sessionManager *scs.SessionManager, repository *repository.Queries, kubeClient *kube.KubeClient) *Handler {
	return &Handler{
		AppConfig:      appConfig,
		SessionManager: sessionManager,
		Repository:     repository,
		KubeClient:     kubeClient,
	}
}
