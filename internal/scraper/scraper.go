package scraper

import (
	"context"
	"log"

	"github.com/vmamchur/joblin-scraper/db/generated"
	"github.com/vmamchur/joblin-scraper/internal/broadcaster"
)

type Scraper interface {
	Name() string
	Scrape() ([]generated.CreateVacancyParams, error)
}

type ScraperService struct {
	scrapers    []Scraper
	broadcaster broadcaster.Broadcaster
}

func NewScraper(broadcaster broadcaster.Broadcaster, djEmail string, djPassword string) *ScraperService {
	return &ScraperService{
		broadcaster: broadcaster,
		scrapers: []Scraper{
			DjinniScraper{email: djEmail, password: djPassword},
		},
	}
}

func (s *ScraperService) Run(q *generated.Queries) {
	for _, scr := range s.scrapers {
		scrName := scr.Name()
		log.Printf("[%s] Starting scrape", scrName)

		vacancies, err := scr.Scrape()
		if err != nil {
			log.Printf("[%s] Scrape failed: %v\n", scrName, err)
			continue
		}
		log.Printf("[%s] Scraped %d vacancies", scrName, len(vacancies))

		for _, v := range vacancies {
			_, err := q.CreateVacancy(context.Background(), v)
			if err != nil {
				log.Printf("[%s] Failed to save: %s - %v\n", scrName, v.Url, err)
				continue
			}
			log.Printf("[%s] Saved: %s\n", scrName, v.Url)

			err = s.broadcaster.Broadcast(v)
			if err != nil {
				log.Printf("[%s] Failed to broadcast: %s - %v\n", scrName, v.Url, err)
				continue
			}
			log.Printf("[%s] Broadcasted: %s\n", scrName, v.Url)
		}

		log.Printf("[%s] Scraping completed\n", scrName)
	}
}
