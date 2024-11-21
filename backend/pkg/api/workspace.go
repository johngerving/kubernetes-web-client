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

// postWorkspaceHandler creates a new workspace for
// the user.
func (s *Server) postWorkspaceHandler(c *gin.Context) {
	// Get user id
	userId := c.MustGet("user").(int32)

	// Bind request body to workspace struct
	workspaceParams := postWorkspaceForm{}
	if err := c.ShouldBind(&workspaceParams); err != nil || !isWorkspaceParamsValid(workspaceParams) {
		log.Printf("invalid response body: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	// Add workspace to db
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

// getWorkspacesHandler gets a list of a user's
// workspaces.
func (s *Server) getWorkspacesHandler(c *gin.Context) {
	userId := c.MustGet("user").(int32)

	workspaces, err := s.repository.ListUserWorkspaces(context.Background(), userId)
	if err != nil {
		log.Printf("error retrieving workspaces: %v\n", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error retrieving workspacse"})
		return
	}

	log.Println(userId)
	log.Println(workspaces)

	c.IndentedJSON(http.StatusOK, workspaces)
}

func isWorkspaceParamsValid(params postWorkspaceForm) bool {
	return params.Name != ""
}
