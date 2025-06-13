package types

import "time"

type Job struct {
	ID      string    `json:"id"`
	Company string    `json:"company"`
	Created time.Time `json:"created"`
}

type HunterResponse struct {
	Data struct {
		Emails []struct {
			Value string `json:"value"`
		} `json:"emails"`
	} `json:"data"`
}
