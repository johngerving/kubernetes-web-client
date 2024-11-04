package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

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
