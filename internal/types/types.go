package types

import "time"

type Job struct {
	ID      string    `json:"job_id"`
	Company string    `json:"company_name"`
	Created time.Time `json:"created"`
}

type HunterResponse struct {
	Data struct {
		Domain       string `json:"domain"`
		Organization string `json:"organization"`
		Description  string `json:"description"`
		Country      string `json:"country"`
		City         string `json:"city"`
		Emails       []struct {
			Value      string `json:"value"`
			FirstName  string `json:"first_name"`
			LastName   string `json:"last_name"`
			Position   string `json:"position"`
			Department string `json:"department"`
			LinkedIn   string `json:"linkedin"`
		} `json:"emails"`
	} `json:"data"`
}

type Company struct {
	ID           int32  `json:"id"`
	Domain       string `json:"domain"`
	Organization string `json:"organization,omitempty"`
	Description  string `json:"description,omitempty"`
	Country      string `json:"country,omitempty"`
	City         string `json:"city,omitempty"`
}

type Contact struct {
	ID         int     `json:"id"`
	Email      string  `json:"email"`
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	Position   *string `json:"position,omitempty"`
	Department *string `json:"department,omitempty"`
	LinkedIn   *string `json:"linkedin,omitempty"`
}

type JobDetails struct {
	JobID     string    `json:"job_id"`
	JobDomain string    `json:"job_domain"`
	Company   Company   `json:"company"`
	Contacts  []Contact `json:"contacts"`
}
