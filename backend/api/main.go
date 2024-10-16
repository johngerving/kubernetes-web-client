package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/johngerving/kubernetes-web-client/backend/api/config"
	"github.com/johngerving/kubernetes-web-client/backend/api/database/db"
	"github.com/johngerving/kubernetes-web-client/backend/api/handler"
	"github.com/johngerving/kubernetes-web-client/backend/api/session"
)

func main() {
	appConfig, err := config.NewConfigFromEnv()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	if appConfig.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create context that listens for interrupt signal from OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Create a Gin router
	router := gin.Default()

	// Initialize database connection
	pool, err := pgxpool.New(context.Background(), appConfig.DBUrl)
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}
	defer pool.Close() // Close connection when done

	sessionManager := session.NewStore(pool)
	queries := db.New(pool)
	fmt.Println(queries)

	h := handler.NewHandler(appConfig, sessionManager)

	router.GET("/auth", h.Auth)
	router.GET("/auth/callback", h.AuthCallback)
	router.GET("/user", h.User)

	// Create an HTTP server listening on the port provided in environment variables
	// using the router we defined
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.Port),
		Handler: sessionManager.LoadAndSave(router),
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
