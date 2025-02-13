package tests

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/database/repository"
	"github.com/johngerving/kubernetes-web-client/backend/pkg/session"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestGETHealth(t *testing.T) {
	t.Run("returns API health status", func(t *testing.T) {
		godotenv.Load(".backend.env")

		apiUrl := os.Getenv("API_URL")
		resp, err := http.Get(apiUrl + "/health")

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
	// tests := struct{}

	godotenv.Load(".backend.env")
	apiUrl := os.Getenv("API_URL")

	for _, tt := range tests {
		t.Run(tt+" returns unauthorized", func(t *testing.T) {
			// Make a request to the user endpoint, with the session cookie in the request
			r, err := http.NewRequest("GET", apiUrl+tt, nil)

			if err != nil {
				t.Fatalf("Error creating %v request: %v", tt, err)
			}

			res, err := http.DefaultClient.Do(r)
			if err != nil {
				t.Fatalf("Error executing %v request: %v", tt, err)
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

		apiUrl := os.Getenv("API_URL")

		client, err := loginUser(pool, apiUrl, "test@example.com")
		if err != nil {
			t.Fatal(err)
		}

		// Make a GET request to the /user endpoint and populate a User struct
		var haveUser repository.User
		statusCode, err := doJSONRequest(client, "GET", apiUrl+"/user", "", &haveUser)
		if err != nil {
			t.Fatal(err)
		}

		require.Equal(t, http.StatusOK, statusCode, "User data response status code should be 200.")
		require.Equal(t, "test@example.com", haveUser.Email, "User data response 'email' field should be 'test@example.com'.")

		// Clear the table once done
		_, err = pool.Exec(context.Background(), "TRUNCATE TABLE sessions, users, workspaces CASCADE")
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

		apiUrl := os.Getenv("API_URL")

		client, err := loginUser(pool, apiUrl, "test@example.com")
		if err != nil {
			t.Fatal(err)
		}

		// Make a request to the user/workspaces endpoint, with the session cookie in the request
		var haveWorkspace repository.Workspace
		statusCode, err := doJSONRequest(client, "POST", apiUrl+"/user/workspaces", `{"name": "test"}`, &haveWorkspace)
		if err != nil {
			t.Fatal(err)
		}

		require.Equal(t, http.StatusOK, statusCode, "Workspace data response status code should be 200.")
		require.Equal(t, "test", haveWorkspace.Name, "User data response 'email' field should be 'test'.")

		// Check if workspace exists in database
		var idFromDB int32
		err = pool.QueryRow(context.Background(), "SELECT id FROM workspaces WHERE id=$1", haveWorkspace.ID).Scan(&idFromDB)
		if err == pgx.ErrNoRows {
			t.Fatalf("Workspace with id %v should exist in the database", haveWorkspace.ID)
		}
		if err != nil {
			t.Fatal(err)
		}

		// Clear the table once done
		_, err = pool.Exec(context.Background(), "TRUNCATE TABLE sessions, users, workspaces CASCADE")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestGETWorkspaces(t *testing.T) {
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

		// Clear database
		_, err = pool.Exec(context.Background(), "TRUNCATE TABLE sessions, users, workspaces CASCADE")
		if err != nil {
			t.Fatal(err)
		}
		apiUrl := os.Getenv("API_URL")

		// Create list of clients and log them in
		clientNames := []string{"test1@example.com", "test2@example.com", "test3@example.com"}
		var clients []*resty.Client

		for _, name := range clientNames {
			hc, err := loginUser(pool, apiUrl, name)
			if err != nil {
				t.Fatal(err)
			}

			client := resty.NewWithClient(hc).
				SetHeader("Content-Type", "application/json")
			clients = append(clients, client)
		}

		// Make requests to the user/workspaces endpoint to populate the users' workspaces
		resp, err := clients[0].R().
			SetBody(`{"name": "user1workspace"}`).
			Post(apiUrl + "/user/workspaces")
		require.Equal(t, nil, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())

		_, err = clients[1].R().
			SetBody(`{"name": "user2workspace1"}`).
			Post(apiUrl + "/user/workspaces")
		require.Equal(t, nil, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())

		_, err = clients[1].R().
			SetBody(`{"name": "user2workspace2"}`).
			Post(apiUrl + "/user/workspaces")
		require.Equal(t, nil, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())

		/******* Test Get Requests *******/
		var body []map[string]string

		resp, err = clients[0].R().
			Get(apiUrl + "/user/workspaces")
		require.Equal(t, nil, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())

		json.Unmarshal(resp.Body(), &body)
		require.Equal(t, 1, len(body), "User 1 should have 1 workspace")
		require.Equal(t, "user1workspace", body[0]["name"], "User 1 should be able to access workspace they created")

		resp, err = clients[1].R().
			Get(apiUrl + "/user/workspaces")
		require.Equal(t, nil, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())

		json.Unmarshal(resp.Body(), &body)
		require.Equal(t, 2, len(body), "User 2 should have 2 workspaces")
		require.Equal(t, "user2workspace1", body[0]["name"])
		require.Equal(t, "user2workspace2", body[1]["name"])

		// Clear the table once done
		_, err = pool.Exec(context.Background(), "TRUNCATE TABLE sessions, users, workspaces CASCADE")
		if err != nil {
			t.Fatal(err)
		}
	})
}

// doJSONRequest takes an HTTP client, a request method, a URL, a request body, and a pointer
// to a responseBody struct. It executes the request and writes the response to the responseBody
// struct. It returns the status code of the request and an error.
func doJSONRequest[T any](client *http.Client, method string, url string, requestBody string, responseBody *T) (int, error) {
	var r *http.Request
	var err error
	if requestBody != "" {
		r, err = http.NewRequest(method, url, bytes.NewBuffer([]byte(requestBody)))
	} else {
		r, err = http.NewRequest(method, url, nil)
	}

	r.Header.Set("Content-Type", "application/json")

	if err != nil {
		return 0, fmt.Errorf("error creating %v request to %v: %v", method, url, err)
	}

	res, err := client.Do(r)
	if err != nil {
		return 0, fmt.Errorf("error executing %v request to %v: %v", method, url, err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal response
	err = json.Unmarshal([]byte(body), responseBody)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling request response: %v", err)
	}

	return res.StatusCode, nil
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
