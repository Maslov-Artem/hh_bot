package storage

import (
	"context"
	"fmt"
	"hh_bot/models"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetJobID(dbpool *pgxpool.Pool, id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
	SELECT EXISTS (SELECT id FROM job_ads WHERE id=$1)
	`
	var exists bool
	err := dbpool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("Failed to check ID existence: %w", err)
	}

	return exists, nil
}

func SaveJobToDB(dbpool *pgxpool.Pool, job *models.JobAd) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `
	INSERT INTO job_ads (
		id, accept_handicapped, accept_incomplete_resumes, accept_kids, accept_temporary,
		allow_messages, alternate_url, apply_alternate_url, approved, archived, area,
		billing_type, code, contacts, department, description, driver_license_types,
		employer, employment_form, experience, fly_in_fly_out_duration, has_test,
		initial_created_at, insider_interview, internship, key_skills, languages, name,
		negotiations_url, night_shifts, premium, professional_roles, published_at,
		relations, response_letter_required, response_url, salary, suitable_resumes_url,
		test, type, video_vacancy, work_format, work_schedule_by_days, working_hours,
		address ) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
		$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32,
		$33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45 )
		ON CONFLICT (id) DO NOTHING;
	`

	_, err := dbpool.Exec(ctx, query,
		job.ID, job.AcceptHandicapped, job.AcceptIncompleteResumes, job.AcceptKids, job.AcceptTemporary,
		job.AllowMessages, job.AlternateURL, job.ApplyAlternateURL, job.Approved, job.Archived, job.Area,
		job.BillingType, job.Code, job.Contacts, job.Department, job.Descrtiption, job.DriverLicenseTypes,
		job.Employer, job.EmploymentForm, job.Experience, job.FlyInFlyOutDuration, job.HasTest,
		job.InitialCreatedAt, job.InsiderInterview, job.Internship, job.KeySkills, job.Languages, job.Name,
		job.NegotiationsUrl, job.NightShifts, job.Premium, job.ProfessionalRoles, job.PublishedAt,
		job.Relations, job.ResponseLetterRequired, job.ResponseURL, job.Salary, job.SuitableResumesURL,
		job.Test, job.Type, job.VideoVacancy, job.WorkFormat, job.WorkScheduleByDays, job.WorkingHours,
		job.Address,
	)

	return err
}

func SaveUnprocessedJobToDB(dbpool *pgxpool.Pool, id string) error {
	query := `
	INSERT INTO processed_job_ads (job_id) VALUES ($1) ON CONFLICT (job_id) DO NOTHING
	`
	_, err := dbpool.Exec(context.Background(), query, id)

	return err
}

func UpdateProcessedJob(dbpool *pgxpool.Pool, id, cover_letter, thinking string) error {
	query := `
	UPDATE processed_job_ads 
	SET cover_letter = $1, thinking = $2, processed = $3 
	WHERE job_id = $4
	`
	_, err := dbpool.Exec(context.Background(), query, cover_letter, thinking, true, id)

	return err
}

func LoadUnprocessedJobs(dbpool *pgxpool.Pool) ([]models.JobAd, error) {
	query := `
	SELECT id, description from job_ads WHERE id IN (SELECT job_id FROM processed_job_ads WHERE processed=false)
	`
	rows, err := dbpool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to load jobs from data base: %w", err)
	}

	var jobAds []models.JobAd

	for rows.Next() {
		var job models.JobAd
		err := rows.Scan(&job.ID, &job.Descrtiption)
		if err != nil {
			log.Printf("failed to retrieve a job: %v", err)
		}
		jobAds = append(jobAds, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobAds, nil
}
