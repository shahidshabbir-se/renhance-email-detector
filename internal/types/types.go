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
