package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

func TestGETHealth(t *testing.T) {
	t.Run("returns API health status", func(t *testing.T) {
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

		t.Log(data)
	})
}
