package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/database/repository"
)

type postWorkspaceForm struct {
	Name string `json:"name"`
}

// valid checks if a postWorkspaceForm struct is valid. It
// returns a map[string]string containing any problems.
func (f *postWorkspaceForm) valid() (problems map[string]string) {
	problems = make(map[string]string)

	if len(f.Name) <= 2 {
		problems["name"] = "Name must be at least 2 characters long"
	}

	return problems
}

// postWorkspaceHandler creates a new workspace for
// the user.
func (s *Server) postWorkspaceHandler(c *gin.Context) {
	// Get user id
	userId := c.MustGet("user").(int32)

	// Bind request body to workspace struct
	workspaceParams := postWorkspaceForm{}
	c.ShouldBind(&workspaceParams)

	if problems := workspaceParams.valid(); len(problems) > 0 {
		log.Printf("workspace param problems: %v\n", problems)
		c.IndentedJSON(http.StatusBadRequest, problems)
		return
	}

	// Add workspace to db
	workspace, err := s.repository.CreateWorkspace(context.Background(), repository.CreateWorkspaceParams{
		Name:  workspaceParams.Name,
		Owner: userId,
	})
	if err != nil {
		log.Printf("error creating workspace: %v\n", err)

		// Check if row already exists in database
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("workspace named %v already exists", workspaceParams.Name)})
		}
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

	c.IndentedJSON(http.StatusOK, workspaces)
}
