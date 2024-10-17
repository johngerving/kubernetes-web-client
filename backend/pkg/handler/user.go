package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type user struct {
	Email string `json:"email"`
}

func (h *Handler) User(c *gin.Context) {
	// Get user data from the session
	email := h.SessionManager.GetString(c.Request.Context(), "email")

	// If found in the session, return the user data
	if email != "" {
		c.IndentedJSON(http.StatusOK, gin.H{"email": email})
		return
	}

	// Respond with unauthorized
	c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
}
