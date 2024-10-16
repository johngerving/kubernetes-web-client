package handler

import (
	"github.com/alexedwards/scs/v2"
	"github.com/johngerving/kubernetes-web-client/backend/api/config"
)

// Struct to hold app information - config, stores, etc.
type Handler struct {
	AppConfig      config.Config
	SessionManager *scs.SessionManager
}

// NewHandler takes in a config.Config struct instance and returns a pointer
// to a Handler struct instance.
func NewHandler(appConfig config.Config, sessionManager *scs.SessionManager) *Handler {
	return &Handler{
		AppConfig:      appConfig,
		SessionManager: sessionManager,
	}
}
