package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/database/repository"
)

var maxOauthStateCookieAge int = 60 * 60 * 24 * 365 // Set max age for OAuth state to a year

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user data from the session
		userId := int32(s.sessionStore.GetInt(c.Request.Context(), "user"))

		if userId == 0 {
			// Respond with unauthorized
			log.Printf("user ID was %v - unauthorized", userId)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}

		// If found in the session, pass the user data along
		c.Set("user", userId)
		c.Next()
	}
}

// authHandler initiates the OAuth flow
func (s *Server) authHandler(c *gin.Context) {
	// Create oauthState cookie
	oauthState := generateOauthState()

	// Set OAuth state cookie with random value and max age that is valid on all paths of the API domain, HTTP only, and secure
	c.SetCookie("oauthstate", oauthState, maxOauthStateCookieAge, "/", s.config.Domain, true, true)

	// Create auth code URL with the OAuth state
	url := s.oauth.AuthCodeURL(oauthState)

	// Redirect to the OAuth page
	c.Redirect(http.StatusFound, url)
}

// generateOauthState generates a random string 16 bytes long
// and returns it.
func generateOauthState() string {
	// Generate a random string for OAuth state
	b := make([]byte, 16)
	rand.Read(b)

	// base64 encode
	return base64.URLEncoding.EncodeToString(b)
}

// authCallbackHandler receives OAuth state and retrieves user information,
// redirecting to the frontend URL when authenticated.
func (s *Server) authCallbackHandler(c *gin.Context) {
	verifier := s.provider.Verifier(&oidc.Config{ClientID: s.oauth.ClientID})
	// Read oauthState from cookie
	oauthState, _ := c.Cookie("oauthstate")
	// Clear the OAuth cookie no matter what
	c.SetCookie("oauthstate", "", maxOauthStateCookieAge, "/", s.config.Domain, true, true)

	// Redirect if state is invalid
	if c.Request.FormValue("state") != oauthState {
		log.Println("error: invalid OAuth state")
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	oauth2Token, err := s.oauth.Exchange(context.Background(), c.Request.URL.Query().Get("code"))
	if err != nil {
		log.Printf("error retrieving OAuth code: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	// Extract ID token from OAuth token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Printf("error extracting OAuth ID token: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	// Parse and verify ID Token payload
	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		log.Printf("error parsing ID token payload: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	// Extract custom claims
	var claims repository.User
	if err := idToken.Claims(&claims); err != nil {
		log.Printf("error extracting OIDC claims: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	// Check if user already exists in database
	user, err := s.repository.FindUserWithEmail(context.Background(), claims.Email)

	if err == pgx.ErrNoRows {
		// If the user isn't in the database, add them
		err = s.repository.CreateUser(context.Background(), claims.Email)
		if err != nil {
			log.Printf("error adding user to database: %v", err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "unable to add user to database"})
			return
		}
	} else if err != nil {
		// Some other error occurred in getting user information
		log.Printf("error retrieving user information: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "unable to retrieve user information"})
		return
	}

	// Create a new session to store the user information
	log.Printf("Putting user ID %v", user.ID)
	s.sessionStore.Put(c.Request.Context(), "user", int(user.ID))

	// Redirect to app URL
	c.Redirect(http.StatusPermanentRedirect, s.config.FrontendURL)
}
