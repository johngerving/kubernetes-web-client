package tests

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/session"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestGETHealth(t *testing.T) {
	t.Run("returns API health status", func(t *testing.T) {
		godotenv.Load(".backend.env")

		port := os.Getenv("PORT")
		resp, err := http.Get("http://localhost:" + port + "/health")

		if err != nil {
			t.Fatalf("error getting /health response: %v", err)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("error reading /health response body: %v", err)
		}

		bodyString := string(body)

		var data map[string]interface{}
		err = json.Unmarshal([]byte(bodyString), &data)
		if err != nil {
			t.Fatalf("unable to unmarshal request response: %v", err)
		}

		require.Equal(t, data["status"], "up", "The API status should be 'up.'")

		var details map[string]interface{}
		require.IsType(t, data["details"], details, "The 'details' field should be a map[string]interface{}.")
		details = data["details"].(map[string]interface{})

		var database map[string]interface{}
		require.IsType(t, details["database"], database, "The 'details['database']' field should be a map[string]interface{}.")
		database = details["database"].(map[string]interface{})

		var databaseType string
		require.IsType(t, database["status"], databaseType, "The database status should be a string.")

		require.Equal(t, database["status"], "up", "The databse status should be 'up.'")
	})
}

func TestAuthEndpoints(t *testing.T) {
	tests := []string{"/user"}

	godotenv.Load(".backend.env")
	port := os.Getenv("PORT")

	for _, tt := range tests {
		t.Run(tt+" returns unauthorized", func(t *testing.T) {
			// Make a request to the user endpoint, with the session cookie in the request
			r, err := http.NewRequest("GET", "http://localhost:"+port+"/user", nil)

			if err != nil {
				t.Fatalf("Error creating /user request: %v", err)
			}

			res, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Fatalf("Error executing /user request: %v", err)
			}

			require.Equal(t, http.StatusUnauthorized, res.StatusCode, "Response status code should be 401.")

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Error reading response body: %v", err)
			}

			// Unmarshal response
			var data map[string]string
			err = json.Unmarshal([]byte(body), &data)
			if err != nil {
				t.Fatalf("unable to unmarshal request response: %v", err)
			}

			require.Equal(t, "unauthorized", data["message"], "Response 'message' field should be 'unauthorized'.")
		})
	}
}

func TestGETUser(t *testing.T) {
	t.Run("returns user information", func(t *testing.T) {
		// Initialize database connection
		dbUrl := os.Getenv("DB_URL")
		if dbUrl == "" {
			t.Fatalf("Error: Database URL must be specified")
		}
		pool, err := pgxpool.New(context.Background(), dbUrl)
		if err != nil {
			t.Fatalf("Failed to initialize database connection: %v", err)
		}
		defer pool.Close() // Close connection when done

		sessionStore := session.NewStore(pool) // New session store

		// Generate test session token
		tokenBytes := make([]byte, 12)
		_, err = rand.Read(tokenBytes)
		if err != nil {
			t.Fatal(err)
		}
		sessionToken := hex.EncodeToString(tokenBytes)

		expiry := time.Now().Add(24 * time.Hour).UTC() // Set session expiration

		// Encode session data
		sessionData, err := sessionStore.Codec.Encode(expiry, map[string]interface{}{
			"email": "test@example.com",
		})
		if err != nil {
			t.Fatalf("Error encoding data: %v", err)
		}

		// Insert session into database
		_, err = pool.Exec(context.Background(), "INSERT INTO sessions VALUES($1, $2, current_timestamp + interval '1 day')", sessionToken, sessionData)
		if err != nil {
			t.Fatal(err)
		}

		godotenv.Load(".backend.env")
		port := os.Getenv("PORT")

		// Create cookie jar for requests
		jar, err := cookiejar.New(nil)
		if err != nil {
			t.Fatalf("Failed to create cookiejar: %v", err)
		}
		client := &http.Client{
			Jar: jar,
		}

		// Create a cookie with the session token
		cookie := &http.Cookie{
			Name:     sessionStore.Cookie.Name,
			Value:    sessionToken,
			Path:     sessionStore.Cookie.Path,
			Domain:   sessionStore.Cookie.Domain,
			Secure:   sessionStore.Cookie.Secure,
			HttpOnly: sessionStore.Cookie.HttpOnly,
			SameSite: sessionStore.Cookie.SameSite,
		}

		cookie.Expires = time.Unix(expiry.Unix()+1, 0)        // Round up to the nearest second.
		cookie.MaxAge = int(time.Until(expiry).Seconds() + 1) // Round up to the nearest second.

		// Add the cookie to the http client
		u, _ := url.Parse("http://localhost:" + port)
		client.Jar.SetCookies(u, []*http.Cookie{cookie})

		// Make a request to the user endpoint, with the session cookie in the request
		r, err := http.NewRequest("GET", "http://localhost:"+port+"/user", nil)

		if err != nil {
			t.Fatalf("Error creating /user request: %v", err)
		}

		res, err := client.Do(r)
		if err != nil {
			t.Fatalf("Error executing /user request: %v", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
		}

		// Unmarshal response
		var data map[string]string
		err = json.Unmarshal([]byte(body), &data)
		if err != nil {
			t.Fatalf("unable to unmarshal request response: %v", err)
		}

		require.Equal(t, "test@example.com", data["email"], "User data response 'email' field should be 'test@example.com'.")

		// Clear the table once done
		_, err = pool.Exec(context.Background(), "TRUNCATE TABLE sessions")
		if err != nil {
			t.Fatal(err)
		}
	})
}
