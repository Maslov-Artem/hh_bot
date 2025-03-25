package main

import (
	"context"
	"flag"
	"fmt"
	"hh_bot/config"
	"hh_bot/jobfetcher"
	"hh_bot/models"
	"hh_bot/processor"
	"hh_bot/storage"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	fetch   = flag.Bool("fetch", false, "flag for downloading job ads")
	process = flag.Bool("process", false, "flag for processing job ads")
)

func main() {
	flag.Parse()
	conf, client, dbpool, err := initialize()

	if err != nil {
		log.Fatalf("Initialization failed: %s", err)
	}
	defer dbpool.Close()
	fmt.Printf("%v, %v\n", *fetch, *process)

	jobQueries := []string{
		"ML Engineer",
		"Data science",
		"Data Scientist",
		"Дата сайентист",
		"Датасайентист",
		"ML",
		"Machine Learning Engineer",
		"ML-инженер",
	}

	if *fetch {
		fmt.Printf("Fetching jobs\n")
		fetchJobAds(client, dbpool, jobQueries, conf.JobAPIURL, conf.JobAPIKey)
	}

	if *process {
		fmt.Printf("Processing jobs\n")
		unprocessedJob, err := storage.LoadUnprocessedJobs(dbpool)
		if err != nil {
			log.Fatal("failed to load unprocessed jobs", err)
		}

		for _, job := range unprocessedJob {
			processAndSaveJob(dbpool, &job, client, conf.Model, conf.SystemPrompt, conf.LLMAPIKey, conf.LLMAPIURL)
			log.Printf("Jobs %s with id %s processed successfully", job.Name, job.ID)

		}
	}
}

func fetchJobAds(client *http.Client, dbpool *pgxpool.Pool, jobQueries []string, queryURL, jobApiKey string) {
	jobMap := make(map[string]int)

	for _, query := range jobQueries {
		jobMap[query] = 0
	}

	for job := range jobMap {
		params := url.Values{}
		params.Add("experience", "noExperience")
		params.Add("experience", "between1And3")
		params.Add("text", job)
		jobsUrl := queryURL + "?" + params.Encode()

		jobs, err := jobfetcher.FetchJobs(client, jobsUrl, jobApiKey)

		if err != nil {
			fmt.Printf("Failed to fetch jobs for %s: %v\n", job, err)
			continue
		}
		jobMap[job] = jobs.Found
		fmt.Printf("Found %d jobs for '%s' query.\n", jobs.Found, job)
	}

	for job := range jobMap {
		for i := 1; i*100 <= jobMap[job]; i++ {
			page := strconv.Itoa(i)
			params := url.Values{}
			params.Add("experience", "noExperience")
			params.Add("experience", "between1And3")
			params.Add("text", job)
			params.Add("per_page", "100")
			params.Add("page", page)

			fetchURL := queryURL + "?" + params.Encode()

			err := fetchAndSaveJobAds(client, fetchURL, jobApiKey, dbpool)

			if err != nil {
				fmt.Printf("Job processing for %s failed: %v\n", job, err)
			}

		}
	}

}

func initialize() (*config.Config, *http.Client, *pgxpool.Pool, error) {
	conf := config.LoadConfig()

	client := &http.Client{Timeout: 20 * time.Second}

	dbpool, err := pgxpool.New(context.Background(), conf.DatabaseURL)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Unable to create connection pool: %w\n", err)
	}

	return conf, client, dbpool, nil
}

func fetchAndSaveJobAds(client *http.Client, fetchURL, jobAPIKey string, dbpool *pgxpool.Pool) error {

	fetchedJobs, err := jobfetcher.FetchJobs(client, fetchURL, jobAPIKey)
	if err != nil {
		return err
	}

	for _, job := range fetchedJobs.Items {
		exist, err := storage.GetJobID(dbpool, job.ID)
		if err != nil {
			log.Printf("failed to check duplicate: %v", err)
		}
		if exist {
			continue
		}
		jobData, err := jobfetcher.ExtractJobData(client, jobAPIKey, job)
		if err != nil {
			log.Printf("failed to extract job data: %v", err)
			continue
		}

		jobData.Descrtiption = processor.RemoveHTMLTags(jobData.Descrtiption)

		err = storage.SaveJobToDB(dbpool, &jobData)
		if err != nil {
			log.Printf("failed to save job to data base: %v", err)
		}

		err = storage.SaveUnprocessedJobToDB(dbpool, job.ID)
		if err != nil {
			log.Printf("failed to save unprocessed job to data base: %v", err)
		}
	}

	return nil

}

func processAndSaveJob(dbpool *pgxpool.Pool, job *models.JobAd, client *http.Client, model, prompt, llmApiKey, llmApiURL string) error {

	description, err := processor.ProcessJob(job, client, model, prompt, llmApiKey, llmApiURL)

	if err != nil {
		return fmt.Errorf("Failed to process job %s: %v\n", job.ID, err)
	} else {
		fmt.Printf("Succsessfully procesed job %s.\n", job.Name)
	}

	err = storage.UpdateProcessedJob(dbpool, job.ID, description[1], description[0])
	if err != nil {
		return fmt.Errorf("Failed to save job %s: %v\n", job.ID, err)
	} else {
		log.Printf("Succsessfully saved processed job %s\n", job.Name)
	}

	return nil
}
