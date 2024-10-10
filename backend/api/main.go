package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	if cfg.env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run()
}
