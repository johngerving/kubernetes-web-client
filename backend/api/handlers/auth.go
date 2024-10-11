package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"

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
	c.Redirect(http.StatusTemporaryRedirect, url)
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
	// Read oauthState from cookie
	oauthState, _ := c.Cookie("oauthstate")

	// Redirect if state is invalid
	if c.Request.FormValue("state") != oauthState {
		log.Println("invalid OAuth state")
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	// Get the user data, redirect if error occurred
	data, err := getUserData(c.Request.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/auth")
		return
	}

	fmt.Printf("userinfo: %s\n", data)
}

func getUserData(code string) ([]byte, error) {
	// Get access token from OAuth Exchange
	token, err := config.AppConfig.OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	// Get userinfo with access token
	response, err := http.Get(config.AppConfig.OAuthUrl + "/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	// Read response and return it
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
