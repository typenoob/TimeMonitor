package utils

import "time"

type Record struct {
	ID         string    `json:"id,omitempty"`
	IPAddress  string    `json:"ip_address,omitempty"`
	LastOkTime time.Time `json:"last_ok_time,omitempty"`
}
