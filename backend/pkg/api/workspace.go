package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
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

// deleteWorkspaceHandler deletes a workspace with a
// given ID.
func (s *Server) deleteWorkspaceHandler(c *gin.Context) {
	userId := c.MustGet("user").(int32)

	// Get the workspace ID
	idParam := c.Param("id")
	workspaceId, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("error in id param: %v\n", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid ID param"})
		return
	}

	// Create params for database delete
	params := repository.DeleteWorkspaceWithIdParams{
		Owner: userId,
		ID:    int32(workspaceId),
	}

	_, err = s.repository.DeleteWorkspaceWithId(context.Background(), params)
	if err == pgx.ErrNoRows {
		log.Printf("row with ID %v owned by user with ID %v does not exist: %v", workspaceId, userId, err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "workspace not found"})
		return
	}
	if err != nil {
		log.Printf("error deleting workspace with ID %v: %v\n", workspaceId, err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error removing workspace"})
		return
	}

	c.Status(http.StatusOK)
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
