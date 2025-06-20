package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/vmamchur/joblin-scraper/db/generated"
)

type Scraper interface {
	Name() string
	Scrape() ([]generated.CreateVacancyParams, error)
}

type ScraperService struct {
	scrapers []Scraper
}

func NewScraper(djEmail string, djPassword string) *ScraperService {
	return &ScraperService{
		scrapers: []Scraper{
			DjinniScraper{email: djEmail, password: djPassword},
		},
	}
}

func (s *ScraperService) Run(q *generated.Queries) {
	for _, scr := range s.scrapers {
		scraperName := scr.Name()
		log.Printf("[%s] Starting scrape", scraperName)

		vacancies, err := scr.Scrape()
		if err != nil {
			log.Printf("[%s] Scrape failed: %v\n", scraperName, err)
			continue
		}
		log.Printf("[%s] Scraped %d vacancies", scraperName, len(vacancies))

		for _, v := range vacancies {
			_, err := q.CreateVacancy(context.Background(), v)
			if err != nil {
				log.Printf("[%s] Failed to save: %s - %v\n", scraperName, v.Url, err)
				continue
			}
			log.Printf("[%s] Saved: %s\n", scraperName, v.Url)

			payload := map[string]string{
				"title":        v.Title,
				"company_name": v.CompanyName,
				"url":          v.Url,
			}

			data, err := json.Marshal(payload)
			if err != nil {
				log.Printf("[%s] Failed to marshal broadcast payload: %v", scraperName, err)
				continue
			}

			resp, err := http.Post("https://43f9-5-58-64-59.ngrok-free.app/broadcast", "application/json", bytes.NewBuffer(data))
			if err != nil {
				log.Printf("[%s] Failed to broadcast vacancy: %s - %v", scraperName, v.Url, err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Printf("[%s] Broadcast returned non-200 status for: %s - %s", scraperName, v.Url, resp.Status)
			} else {
				log.Printf("[%s] Broadcasted: %s", scraperName, v.Url)
			}
		}

		log.Printf("[%s] Scraping completed\n", scraperName)
	}
}
