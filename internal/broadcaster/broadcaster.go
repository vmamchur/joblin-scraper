package broadcaster

import "github.com/vmamchur/joblin-scraper/db/generated"

type Broadcaster interface {
	Broadcast(v generated.CreateVacancyParams) error
}
