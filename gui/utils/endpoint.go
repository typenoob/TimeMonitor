package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	Port = 1234
)

func GetAllRecord(NTPServer string) ([]Record, error) {
	var records []Record
	var err error
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/records", NTPServer, Port))
	if err != nil {
		return records, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return records, err
	}
	if b != nil {
		json.Unmarshal(b, &records)
	}
	return records, nil
}
