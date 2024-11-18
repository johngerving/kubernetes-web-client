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
	userId := c.MustGet("user").(int32)

	user, err := s.repository.FindUserWithId(context.Background(), userId)
	if err == pgx.ErrNoRows {
		log.Printf("error: user with ID %v does not exist in database", userId)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	if err != nil {
		log.Printf("error retrieving user with ID %v from database: %v", userId, err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error retrieving user"})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}
