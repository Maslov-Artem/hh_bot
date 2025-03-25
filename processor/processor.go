package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hh_bot/models"
	"io"

	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const MaxRetries = 3

func RemoveHTMLTags(input string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}

func reThink(input string) ([]string, error) {
	reThink := regexp.MustCompile(`<think>(?s)(.*?)</think>`)
	think := reThink.FindStringSubmatch(input)

	reAfterThink := regexp.MustCompile(`</think>(?s)(.*)`)
	afterThink := reAfterThink.FindStringSubmatch(input)

	if len(think) < 1 || len(afterThink) < 1 {
		return nil, errors.New("Failed to extract cover letter from LLM response")
	}

	return []string{strings.TrimSpace(think[1]), strings.TrimSpace(afterThink[1])}, nil
}

func ProcessJobDesctription(text string, client *http.Client, model, prompt, llmApiKey, llmApiURL string) (string, error) {

	requestPayload, err := createRequestPayload(model, prompt, text)
	if err != nil {
		return "", fmt.Errorf("failed to create request payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body []byte
	var lastErr error

	for attempt := range MaxRetries {

		resp, err := makeGroqApiCall(client, ctx, llmApiKey, llmApiURL, requestPayload)
		if err != nil {
			lastErr = fmt.Errorf("failed to make API call: %w", err)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			apiResp, err := decodeApiResponse(resp)
			if err != nil {
				lastErr = fmt.Errorf("Failed to decode API response %w", err)
				continue
			}
			return apiResp, nil
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			retryTimeHeader := resp.Header.Get("retry-after")
			retryTime, err := strconv.Atoi(retryTimeHeader)
			if err != nil {
				lastErr = fmt.Errorf("failed to converet retry time header: %w", err)
				continue
			}

			fmt.Printf("Rate limited - retrying request in %d seconds (attempt %d/%d)\n",
				retryTime, attempt+1, MaxRetries)

			select {
			case <-time.After(time.Duration(retryTime) * time.Second):
				continue
			case <-ctx.Done():
				return "", ctx.Err()
			}

		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %v", err)
			continue
		}

		lastErr = fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, body)
	}

	return "", fmt.Errorf("max retries (%d) exceeded, last error: %w", MaxRetries, lastErr)
}

func decodeApiResponse(response *http.Response) (string, error) {
	var apiResponse models.GroqAPIResponse
	err := json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices found in response")
	}

	return apiResponse.Choices[0].Message.Content, nil
}

func createRequestPayload(model, system, text string) ([]byte, error) {
	requestPayload := models.GroqAPIRequest{
		Messages: []models.Message{
			{
				Role:    "system",
				Content: system,
			},
			{
				Role:    "user",
				Content: text,
			},
		},
		Model: model,
	}

	payloadBytes, err := json.Marshal(requestPayload)

	if err != nil {
		return nil, err
	}

	return payloadBytes, nil
}

func makeGroqApiCall(client *http.Client, ctx context.Context, apiKey, apiURL string, requestPayload []byte) (*http.Response, error) {

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create a request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ProcessJob(job *models.JobAd, client *http.Client, model, prompt, llmApiKey, llmApiURL string) ([]string, error) {

	processedText, err := ProcessJobDesctription(job.Descrtiption, client, model, prompt, llmApiKey, llmApiURL)
	if err != nil {
		return nil, err
	}

	text, err := reThink(processedText)
	if err != nil {
		fmt.Printf("Failed extact data from LLM: %s\n", err)
		return nil, err
	}

	return text, nil
}
