package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/database/repository"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/kube"
	"golang.org/x/oauth2"
)

type Server struct {
	router       *gin.Engine         // Gin router
	config       *Config             // General app config
	oauth        *oauth2.Config      // OAuth config
	provider     *oidc.Provider      // OIDC provider
	sessionStore *scs.SessionManager // Session store
	repository   *repository.Queries // Database
	kube         *kube.Client        // Kubernetes client
}

// NewServer takes a Config, oauth2.Config, oidc.Provider, scs.SessionManager, repository.Queries, and kube.Client
// and returns a Server.
func NewServer(config *Config, oauth *oauth2.Config, provider *oidc.Provider, sessionStore *scs.SessionManager, repo *repository.Queries, kube *kube.Client) (*Server, error) {
	srv := &Server{
		router:       gin.Default(),
		config:       config,
		oauth:        oauth,
		provider:     provider,
		sessionStore: sessionStore,
		repository:   repo,
		kube:         kube,
	}

	return srv, nil
}

// registerHandlers registers the routes of the API.
func (s *Server) registerHandlers() {
	s.router.GET("/auth", s.authHandler)
	s.router.GET("/auth/callback", s.authCallbackHandler)
	s.router.GET("/user", s.userHandler)
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() *http.Server {
	s.registerHandlers()

	// Create context that listens for interrupt signal from OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Create an HTTP server listening on the port provided in environment variables
	// using the router we defined
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.sessionStore.LoadAndSave(s.router),
	}

	// Initialize the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for interrupt signal
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown
	stop()
	log.Println("Shutting down server gracefully")

	// Context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

	return srv
}
