package models

import (
	"encoding/json"
	"time"
)

type Address struct {
	Building      string         `json:"building"`
	City          string         `json:"city"`
	Description   string         `json:"description"`
	Lat           float32        `json:"lat"`
	Lng           float32        `json:"lng"`
	MetroStations []MetroStation `json:"metro_stations"`
	Street        string         `json:"street"`
}

type Area struct {
	NamedEntity
	URL string `json:"url"`
}

type Contacts struct {
	CallTrackingEnabled bool     `json:"call_tracking_enabled"`
	Email               *string  `json:"email"`
	Name                *string  `json:"name"`
	Phones              *[]Phone `json:"phones"`
}

type Counter struct {
	Responses      int `json:"responses"`
	TotalResponses int `json:"total_responses"`
}

type DriverLicenseType struct {
	ID string `json:"id"`
}

type Employer struct {
	NamedEntity
	Accredited_it_employer bool `json:"accredited_it_employer"`
	EmployerRating         struct {
		ReviewsCount interface{} `json:"reviews_count"`
		TotalRating  string      `json:"total_rating"`
	} `json:"employer_rating"`
	AlternateURL *string `json:"alternate_url"`
	LogoURLs     struct {
		Size240  string `json:"240"`
		Original string `json:"original"`
		Size90   string `json:"90"`
	} `json:"logo_urls"`
	Trusted           bool   `json:"trusted"`
	URL               string `json:"url"`
	VacanciesURL      string `json:"vacancies_url"`
	Blacklisted       bool   `json:"blacklisted"`
	ApplicantServices struct {
		TargetEmployer struct {
			Count int `json:"count"`
		} `json:"target_employer"`
	} `json:"applicant_services"`
}

type InsiderInterview struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type KeySkill struct {
	Name string `json:"name"`
}

type Language struct {
	NamedEntity
	Level NamedEntity `json:"level"`
}

type MetroStation struct {
	Lat         float32 `json:"lat"`
	LineID      string  `json:"line_id"`
	LineName    string  `json:"line_name"`
	Lng         float32 `json:"lng"`
	StationID   string  `json:"station_id"`
	StationName string  `json:"station_name"`
}

type NamedEntity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Phone struct {
	City      string  `json:"city"`
	Comment   *string `json:"comment"`
	Country   string  `json:"country"`
	Formatted string  `json:"formatted"`
	Number    string  `json:"number"`
}

type Salary struct {
	From     *int   `json:"from"`
	To       *int   `json:"to"`
	Currency string `json:"currency"`
	Gross    bool   `json:"gross"`
}

type Snippet struct {
	Requirement    string `json:"requirement"`
	Responsibility string `json:"responsibility"`
}

type Test struct {
	ID       *string
	Required bool
}

type VideoVacancy struct {
	CoverPicture struct {
		ResizedHeight int    `json:"resized_height"`
		ResizedPath   string `json:"resized_path"`
		ResizedWidth  int    `json:"resized_width"`
	} `json:"cover_picture"`
	SnippetPictureURL *string `json:"snippet_picture_url"`
	SnippetVideoURL   *string `json:"snippet_video_url"`
	VideoURL          string  `json:"video_url"`
}

type JobAd struct {
	AcceptHandicapped       bool                `json:"accept_handicapped"`
	AcceptIncompleteResumes bool                `json:"accept_incomplete_resumes"`
	AcceptKids              bool                `json:"accept_kids"`
	AcceptTemporary         bool                `json:"accept_temporary"`
	AllowMessages           bool                `json:"allow_messages"`
	AlternateURL            string              `json:"alternate_url"`
	ApplyAlternateURL       string              `json:"apply_alternate_url"`
	Approved                bool                `json:"approved"`
	Archived                bool                `json:"archived"`
	Area                    Area                `json:"area"`
	BillingType             NamedEntity         `json:"billing_type"`
	Code                    *string             `json:"code"`
	Contacts                *Contacts           `json:"contacts"`
	Department              *NamedEntity        `json:"department"`
	Descrtiption            string              `json:"description"`
	DriverLicenseTypes      []DriverLicenseType `json:"driver_license_types"`
	Employer                Employer            `json:"employer"`
	EmploymentForm          NamedEntity         `json:"employment_form"`
	Experience              NamedEntity         `json:"experience"`
	FlyInFlyOutDuration     []NamedEntity       `json:"fly_in_fly_out_duration"`
	HasTest                 bool                `json:"has_test"`
	ID                      string              `json:"id"`
	InitialCreatedAt        CustomTime          `json:"initial_created_at"`
	InsiderInterview        InsiderInterview    `json:"insider_interview"`
	Internship              bool                `json:"internship"`
	KeySkills               []KeySkill          `json:"key_skills"`
	Languages               []Language          `json:"languages"`
	Name                    string              `json:"name"`
	NegotiationsUrl         *string             `json:"negotiations_url"`
	NightShifts             bool                `json:"night_shifts"`
	Premium                 bool                `json:"premium"`
	ProfessionalRoles       []NamedEntity       `json:"professional_roles"`
	PublishedAt             CustomTime          `json:"published_at"`
	Relations               *[]string           `json:"relations"`
	ResponseLetterRequired  bool                `json:"response_letter_required"`
	ResponseURL             *string             `json:"response_url"`
	Salary                  Salary              `json:"salary"`
	SuitableResumesURL      *string             `json:"suitable_resumes_url"`
	Test                    Test                `json:"test"`
	Type                    NamedEntity         `json:"type"`
	VideoVacancy            VideoVacancy        `json:"video_vacancy"`
	WorkFormat              []NamedEntity       `json:"work_format"`
	WorkScheduleByDays      []NamedEntity       `json:"work_schedule_by_days"`
	WorkingHours            []NamedEntity       `json:"working_hours"`
	Address                 Address             `json:"address"`
}

type JobListing struct {
	AcceptIncompleteResumes bool             `json:"accept_incomplete_resumes"`
	AcceptTemporary         bool             `json:"accept_temporary"`
	Address                 Address          `json:"address"`
	AlternateURL            string           `json:"alternate_url"`
	ApplyAlternateURL       string           `json:"apply_alternate_url"`
	Archived                bool             `json:"archived"`
	Area                    Area             `json:"area"`
	Contacts                Contacts         `json:"contacts"`
	CreatedAt               CustomTime       `json:"created_at"`
	Department              NamedEntity      `json:"department"`
	Employer                Employer         `json:"employer"`
	FlyInFlyOutDuration     []NamedEntity    `json:"fly_in_fly_out_duration"`
	HasTest                 bool             `json:"has_test"`
	ID                      string           `json:"id"`
	InsiderInterview        InsiderInterview `json:"insider_interview"`
	Internship              bool             `json:"internship"`
	MetroStations           []MetroStation   `json:"metro_stations"`
	Name                    string           `json:"name"`

	NightShifts            bool          `json:"night_shifts"`
	Premium                bool          `json:"premium"`
	ProfessionalRoles      []NamedEntity `json:"professional_roles"`
	PublishedAt            CustomTime    `json:"published_at"`
	Relations              *[]string     `json:"relations"`
	ResponseLetterRequired bool          `json:"response_letter_required"`
	ResponseURL            *string       `json:"response_url"`
	Salary                 Salary        `json:"salary"`
	SortPointDistance      *int          `json:"sort_point_distance"`
	Type                   NamedEntity   `json:"type"`
	URL                    string        `json:"url"`
	WorkFormat             []NamedEntity `json:"work_format"`
	WorkScheduleByDays     []NamedEntity `json:"work_schedule_by_days"`
	WorkingHours           []NamedEntity `json:"working_hours"`
	Counters               []Counter     `json:"counters"`
	EmploymentForm         NamedEntity   `json:"employment_form"`
	Experience             NamedEntity   `json:"experience"`
	Snippet                Snippet       `json:"snippet"`
	ShowLogoInSearch       *bool         `json:"show_logo_in_search"`
	VideoVacancy           VideoVacancy  `json:"video_vacancy"`
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type JobSearchResponse struct {
	Items   []JobListing `json:"items"`
	Found   int          `json:"found"`
	Page    int          `json:"page"`
	Pages   int          `json:"pages"`
	PerPage int          `json:"per_page"`
}

type GroqAPIRequest struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqAPIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index             int     `json:"index"`
	Message           Message `json:"message"`
	Logprobs          bool    `json:"logprobs"`
	FinishReason      string  `json:"finish_reason"`
	Usage             Usage   `json:"usage"`
	SystemFingerprint string  `json:"system_fingerprint"`
	XGroq             struct {
		ID string `json:"id"`
	} `json:"xgroq"`
}

type Usage struct {
	QueueTime        float64 `json:"queue_time"`
	PromptTokens     int     `json:"prompt_tokens"`
	PromptTime       float64 `json:"prompt_time"`
	CompletionTokens int     `json:"completion_tokens"`
	CompletionTime   float64 `json:"completion_time"`
	TotalTokens      int     `json:"total_tokens"`
	TotalTime        float64 `json:"total_time"`
}

type CustomTime time.Time

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return err
	}
	const inputTimeFormat = "2006-01-02T15:04:05-0700"
	t, err := time.Parse(inputTimeFormat, timeStr)
	if err != nil {
		return err
	}

	*ct = CustomTime(t)
	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(ct).Format(time.RFC3339))
}

func (ct CustomTime) String() string {
	return time.Time(ct).Format(time.RFC3339)
}
