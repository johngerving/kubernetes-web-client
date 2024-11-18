package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// userHandler gets the information of the user based on their session.
func (s *Server) userHandler(c *gin.Context) {
	email := c.MustGet("user").(string)

	user, err := s.repository.FindUserWithEmail(context.Background(), email)
	if err == pgx.ErrNoRows {
		log.Printf("error: user with email %v does not exist in database", email)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	if err != nil {
		log.Printf("error retrieving user with email %v from database: %v", email, err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error retrieving user"})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}
