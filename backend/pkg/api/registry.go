package api

import (
	"github.com/alexliesenfeld/health"
	"github.com/gin-gonic/gin"
)

// Interface for registering handlers
type HandlerRegistry interface {
	RegisterHandlers(s *Server)
}

type MainServerRegistry struct{}

// RegisterHandlers registers the routes of the main server API.
func (r MainServerRegistry) RegisterHandlers(s *Server) {
	unAuthed := s.router.Group("")
	{
		unAuthed.GET("/health", gin.WrapF(health.NewHandler(s.healthChecker))) // Create a handler for a health check and make it an endpoint
		unAuthed.GET("/auth", s.authHandler)
		unAuthed.GET("/auth/callback", s.authCallbackHandler)
	}

	authed := s.router.Group("")
	{
		authed.Use(s.authMiddleware())

		authed.GET("/user", s.userHandler)
	}
}
