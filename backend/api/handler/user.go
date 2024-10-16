package handler

import (
	"log"

	"github.com/gin-gonic/gin"
)

type user struct {
	Email string `json:"email"`
}

func (h *Handler) User(c *gin.Context) {
	// Get user data from the session
	user := h.SessionManager.GetString(c.Request.Context(), "user")

	log.Println(user)
}
