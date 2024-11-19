package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/database/repository"
)

type postWorkspaceForm struct {
	Name string `json:"name" binding:"required"`
}

func (s *Server) postWorkspaceHandler(c *gin.Context) {
	userId := c.MustGet("user").(int32)

	workspaceParams := postWorkspaceForm{}

	if err := c.ShouldBind(&workspaceParams); err != nil || !isWorkspaceParamsValid(workspaceParams) {
		log.Printf("invalid response body: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	workspace, err := s.repository.CreateWorkspace(context.Background(), repository.CreateWorkspaceParams{
		Name:  workspaceParams.Name,
		Owner: userId,
	})
	if err != nil {
		log.Printf("error creating workspace: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error creating workspace"})
		return
	}

	c.IndentedJSON(http.StatusOK, workspace)
}

func isWorkspaceParamsValid(params postWorkspaceForm) bool {
	return params.Name != ""
}
