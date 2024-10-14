package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/johngerving/kubernetes-web-client/backend/api/config"
)

func AuthHandler(c *gin.Context) {
	// Create oauthState cookie
	oauthState := generateOauthState()

	maxAge := 60 * 60 * 24 * 365 // Set max age to a year
	// Set OAuth state cookie with random value and max age that is valid on all paths of the API domain, HTTP only, and secure
	c.SetCookie("oauthstate", oauthState, maxAge, "/", config.AppConfig.Domain, true, true)

	// Create auth code URL with the OAuth state
	url := config.AppConfig.OAuthConfig.AuthCodeURL(oauthState)

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

func AuthCallbackHandler(c *gin.Context) {
	verifier := config.AppConfig.Provider.Verifier(&oidc.Config{ClientID: config.AppConfig.OAuthConfig.ClientID})
	// Read oauthState from cookie
	oauthState, _ := c.Cookie("oauthstate")

	// Redirect if state is invalid
	if c.Request.FormValue("state") != oauthState {
		log.Println("invalid OAuth state")
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	oauth2Token, err := config.AppConfig.OAuthConfig.Exchange(context.Background(), c.Request.URL.Query().Get("code"))
	if err != nil {
		log.Println("error retrieving OAuth code")
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	// Extract ID token from OAuth token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Println("error extracting OAuth ID token")
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	// Parse and verify ID Token payload
	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		log.Println("error parsing ID token payload")
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	// Extract custom claims
	var claims struct {
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		log.Println("error extracting OIDC claims")
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	fmt.Println(claims)
	// Redirect to app URL
	c.Redirect(http.StatusPermanentRedirect, config.AppConfig.AppUrl)
}
