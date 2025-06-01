package scraper

import (
	"log"

	"github.com/vmamchur/joblin-scraper/db/generated"
)

type Scraper interface {
	Scrape() error
}

type ScraperService struct {
	q        *generated.Queries
	scrapers []Scraper
}

func NewScraper(q *generated.Queries, djEmail string, djPassword string) *ScraperService {
	return &ScraperService{
		scrapers: []Scraper{
			DjinniScraper{q: q, email: djEmail, password: djPassword},
		},
	}
}

func (s *ScraperService) Run() {
	for _, scr := range s.scrapers {
		err := scr.Scrape()
		if err != nil {
			log.Printf("Error scraping: %v\n", err)
			continue
		}
	}
}
