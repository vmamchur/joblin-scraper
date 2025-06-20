package broadcaster

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/vmamchur/joblin-scraper/db/generated"
)

type TelegramBroadcaster struct {
	Endpoint string
}

func NewTelegramBroadcaster(endpoint string) *TelegramBroadcaster {
	return &TelegramBroadcaster{endpoint}
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

	resp, err := http.Post(tb.Endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
