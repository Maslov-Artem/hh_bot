package jobfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"hh_bot/models"
	"io"
	"net/http"
	"time"
)

func FetchJobs(client *http.Client, url, jobAPIKey string) (*models.JobSearchResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+jobAPIKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "Aplication aplier")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var searchResponse models.JobSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &searchResponse, nil
}

func ExtractJobData(client *http.Client, jobAPIKey string, job models.JobListing) (models.JobAd, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, job.URL, nil)
	if err != nil {
		return models.JobAd{}, fmt.Errorf("failed to create request for job %s: %v", job.ID, err)
	}
	req.Header.Add("Authorization", "Bearer "+jobAPIKey)
	req.Header.Add("Contetn-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return models.JobAd{}, fmt.Errorf("failed to fetch job %s: %v", job.ID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.JobAd{}, fmt.Errorf("unexpected status code for job %s: %d", job.ID, resp.StatusCode)
	}

	var jobData models.JobAd

	err = json.NewDecoder(resp.Body).Decode(&jobData)
	if err != nil {
		return models.JobAd{}, fmt.Errorf("failed to decode job %s: %v", job.ID, err)
	}

	return jobData, nil
}
