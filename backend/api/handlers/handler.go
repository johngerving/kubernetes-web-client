package handlers

import "github.com/johngerving/kubernetes-web-client/backend/api/config"

// Struct to hold app information - config, stores, etc.
type Handler struct {
	appConfig config.Config
}

// NewHandler takes in a config.Config struct instance and returns a pointer
// to a Handler struct instance.
func NewHandler(appConfig config.Config) *Handler {
	return &Handler{
		appConfig: appConfig,
	}
}
