package processor_test

import (
	"encoding/json"
	"hh_bot/config"
	"hh_bot/models"
	"hh_bot/processor"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMakeGroqApiCall(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mock-endpoint" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		// Simulate rate limiting (429 Too Many Requests)
		if r.Header.Get("Authorization") == "Bearer invalid_key" {
			w.Header().Set("Retry-After", "3")
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Simulating a successful response
		resp := models.GroqAPIResponse{
			Choices: []models.Choice{
				{Message: models.Message{
					Content: "Sucssessfull call",
				},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	client := &http.Client{Timeout: 5 * time.Second}

	// Test case: Successful API call
	t.Run("Successful API call", func(t *testing.T) {
		conf := &config.Config{LLMAPIURL: mockServer.URL + "/mock-endpoint", LLMAPIKey: "valid_key"}
		text := ""
		response, err := processor.ProcessJobDesctription(text, client, conf.Model, conf.SystemPrompt, conf.LLMAPIKey, conf.LLMAPIURL)
		t.Log(response)
		if err != nil {
			t.Fatalf("API call failed: %v", err)
		}

	})

	// Test case: Rate-limited API call (429 Too Many Requests)
	t.Run("Rate-limited API call", func(t *testing.T) {
		conf := &config.Config{LLMAPIKey: "invalid_key", LLMAPIURL: mockServer.URL + "/mock-endpoint"}
		text := ""
		_, err := processor.ProcessJobDesctription(text, client, conf.Model, conf.SystemPrompt, conf.LLMAPIKey, conf.LLMAPIURL)
		if err == nil {
			t.Fatalf("API call failed: %v", err)
		}
	})
}
