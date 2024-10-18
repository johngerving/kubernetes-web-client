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
	"github.com/johngerving/kubernetes-web-client/backend/pkg/config"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/database/repository"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/handler"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/kube"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/session"
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

	kubeClient, err := kube.NewKubeClient(appConfig.KubeConfig)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v\n", err)
	}

	// Initialize database connection
	pool, err := pgxpool.New(context.Background(), appConfig.DBUrl)
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}
	defer pool.Close() // Close connection when done

	sessionManager := session.NewStore(pool)
	repository := repository.New(pool)

	h := handler.NewHandler(appConfig, sessionManager, repository, kubeClient)

	// Create context that listens for interrupt signal from OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Create a Gin router
	router := gin.Default()

	router.GET("/auth", h.Auth)
	router.GET("/auth/callback", h.AuthCallback)
	router.GET("/user", h.User)
	router.GET("/pods", func(c *gin.Context) {
		pods, err := h.KubeClient.ListPods(context.Background())
		if err != nil {
			log.Printf("error listing pods: %v\n", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error listing pods"})
			return
		}
		type Pod struct {
			Name string `json:"name"`
		}
		podList := make([]Pod, len(pods))
		for i := range pods {
			podList[i].Name = pods[i].Name
		}
		c.IndentedJSON(http.StatusOK, podList)
	})

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
