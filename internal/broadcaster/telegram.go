package broadcaster

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/vmamchur/joblin-scraper/db/generated"
)

type TelegramBroadcaster struct {
	ApiUrl string
	ApiKey string
}

func NewTelegramBroadcaster(apiUrl string, apiKey string) *TelegramBroadcaster {
	return &TelegramBroadcaster{ApiUrl: apiUrl, ApiKey: apiKey}
}

func (tb TelegramBroadcaster) Broadcast(v generated.CreateVacancyParams) error {
	payload := map[string]string{
		"title":        v.Title,
		"company_name": v.CompanyName,
		"url":          v.Url,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", tb.ApiUrl+"/broadcast", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", tb.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
