package utils

type Record struct {
	ID         string `json:"id,omitempty"`
	IPAddress  string `json:"ip_address,omitempty"`
	LastOkTime string `json:"last_ok_time,omitempty"`
}
