package tests

import (
	"bytes"
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
	"github.com/johngerving/kubernetes-web-client/backend/pkg/database/repository"
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
	godotenv.Load(".backend.env")
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

		port := os.Getenv("PORT")

		client, err := loginUser(pool, "http://localhost:"+port, "test@example.com")
		if err != nil {
			t.Fatal(err)
		}

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
		var haveUser repository.User
		err = json.Unmarshal([]byte(body), &haveUser)
		if err != nil {
			t.Fatalf("unable to unmarshal request response: %v", err)
		}

		require.Equal(t, http.StatusOK, res.StatusCode, "User data response status code should be 200.")
		require.Equal(t, "test@example.com", haveUser.Email, "User data response 'email' field should be 'test@example.com'.")

		// Clear the table once done
		_, err = pool.Exec(context.Background(), "TRUNCATE TABLE sessions, users CASCADE")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestPOSTWorkspace(t *testing.T) {
	godotenv.Load(".backend.env")
	t.Run("creates a workspace", func(t *testing.T) {
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

		port := os.Getenv("PORT")

		client, err := loginUser(pool, "http://localhost:"+port, "test@example.com")
		if err != nil {
			t.Fatal(err)
		}

		// Make a request to the user/workspaces endpoint, with the session cookie in the request
		var jsonStr = []byte(`{"name": "test"}`)
		r, err := http.NewRequest("POST", "http://localhost:"+port+"/user/workspaces", bytes.NewBuffer(jsonStr))
		r.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Fatalf("Error creating /user/workspaces request: %v", err)
		}

		res, err := client.Do(r)
		if err != nil {
			t.Fatalf("Error executing /user/workspaces request: %v", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
		}

		// Unmarshal response
		var haveWorkspace repository.Workspace
		err = json.Unmarshal([]byte(body), &haveWorkspace)
		if err != nil {
			t.Fatalf("unable to unmarshal request response: %v", err)
		}

		require.Equal(t, http.StatusOK, res.StatusCode, "Workspace data response status code should be 200.")
		require.Equal(t, "test", haveWorkspace.Name, "User data response 'email' field should be 'test'.")

		// Clear the table once done
		_, err = pool.Exec(context.Background(), "TRUNCATE TABLE sessions, users CASCADE")
		if err != nil {
			t.Fatal(err)
		}
	})
}

// loginUser takes a database connection, a domain, and an email and creates and logs in
// a user. It returns an http.Client with a cookie authenticating the user.
func loginUser(pool *pgxpool.Pool, domain string, email string) (*http.Client, error) {
	sessionStore := session.NewStore(pool) // New session store

	// Generate test session token
	tokenBytes := make([]byte, 12)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return nil, err
	}
	sessionToken := hex.EncodeToString(tokenBytes)

	expiry := time.Now().Add(24 * time.Hour).UTC() // Set session expiration

	wantUser := repository.User{
		Email: email,
	}

	var userId int32
	// Insert user into database
	err = pool.QueryRow(context.Background(), "INSERT INTO users (email) VALUES ($1) RETURNING id", wantUser.Email).Scan(&userId)
	if err != nil {
		return nil, err
	}

	// Encode session data
	sessionData, err := sessionStore.Codec.Encode(expiry, map[string]interface{}{
		"user": int(userId),
	})
	if err != nil {
		return nil, err
	}

	// Insert session into database
	_, err = pool.Exec(context.Background(), "INSERT INTO sessions VALUES($1, $2, current_timestamp + interval '1 day')", sessionToken, sessionData)
	if err != nil {
		return nil, err
	}

	// Create cookie jar for requests
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
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
	u, _ := url.Parse(domain)
	client.Jar.SetCookies(u, []*http.Cookie{cookie})

	return client, nil
}
