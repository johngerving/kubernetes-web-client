package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/api"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/controller"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/database/repository"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/oauth"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/session"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Get server config
	serverCfg, err := api.NewConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Set Gin mode to release if in production environment
	if serverCfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Get OAuth config and OIDC provider
	oauth, provider, err := oauth.NewConfigAndProviderFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Get cluster Controller
	controller, err := controller.NewControllerFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database connection
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatalf("Error: Database URL must be specified")
	}
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}
	defer pool.Close() // Close connection when done

	sessionStore := session.NewStore(pool) // New session store
	repository := repository.New(pool)     // New database repository

	// Create the server
	srv, err := api.NewServer(serverCfg, oauth, provider, sessionStore, repository, controller)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	// Create main server registry
	registry := api.MainServerRegistry{}

	// Listen on the server, using the main server registry
	srv.ListenAndServe(registry)
}
