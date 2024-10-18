package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type user struct {
	Email string `json:"email"`
}

// userHandler gets the information of the user based on their session.
func (s *Server) userHandler(c *gin.Context) {
	// Get user data from the session
	email := s.sessionStore.GetString(c.Request.Context(), "email")

	// If found in the session, return the user data
	if email != "" {
		c.IndentedJSON(http.StatusOK, gin.H{"email": email})
		return
	}

	// Respond with unauthorized
	c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
}
