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
	email := c.MustGet("user").(string)

	c.IndentedJSON(http.StatusOK, gin.H{"email": email})
}
