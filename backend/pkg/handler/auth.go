package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var maxOauthStateCookieAge int = 60 * 60 * 24 * 365 // Set max age for OAuth state to a year

func (h *Handler) Auth(c *gin.Context) {
	// Create oauthState cookie
	oauthState := generateOauthState()

	// Set OAuth state cookie with random value and max age that is valid on all paths of the API domain, HTTP only, and secure
	c.SetCookie("oauthstate", oauthState, maxOauthStateCookieAge, "/", h.AppConfig.Domain, true, true)

	// Create auth code URL with the OAuth state
	url := h.AppConfig.OAuthConfig.AuthCodeURL(oauthState)

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

func (h *Handler) AuthCallback(c *gin.Context) {
	verifier := h.AppConfig.Provider.Verifier(&oidc.Config{ClientID: h.AppConfig.OAuthConfig.ClientID})
	// Read oauthState from cookie
	oauthState, _ := c.Cookie("oauthstate")
	// Clear the OAuth cookie no matter what
	c.SetCookie("oauthstate", "", maxOauthStateCookieAge, "/", h.AppConfig.Domain, true, true)

	// Redirect if state is invalid
	if c.Request.FormValue("state") != oauthState {
		log.Println("error: invalid OAuth state")
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	oauth2Token, err := h.AppConfig.OAuthConfig.Exchange(context.Background(), c.Request.URL.Query().Get("code"))
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
	var claims user
	if err := idToken.Claims(&claims); err != nil {
		log.Printf("error extracting OIDC claims: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	// Check if user already exists in database
	user, err := h.Repository.FindUserWithEmail(context.Background(), claims.Email)

	if err == pgx.ErrNoRows {
		// If the user isn't in the database, add them
		err = h.Repository.CreateUser(context.Background(), claims.Email)
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
	h.SessionManager.Put(c.Request.Context(), "email", user.Email)

	// Redirect to app URL
	c.Redirect(http.StatusPermanentRedirect, h.AppConfig.AppUrl)
}
