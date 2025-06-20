package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"github.com/vmamchur/joblin-scraper/config"
	"github.com/vmamchur/joblin-scraper/db/generated"
	"github.com/vmamchur/joblin-scraper/internal/broadcaster"
	"github.com/vmamchur/joblin-scraper/internal/scraper"
)

func main() {
	cfg := config.Load()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Password,
		cfg.DB.Host, cfg.DB.Port,
		cfg.DB.Name, cfg.DB.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %s\n", err)
	}
	defer db.Close()

	q := generated.New(db)
	tgBroadcaster := broadcaster.NewTelegramBroadcaster(cfg.TgBotUrl)

	scraper := scraper.NewScraper(tgBroadcaster, cfg.Djinni.Email, cfg.Djinni.Password)

	log.Println("Scraper scheduler started")
	scraper.Run(q)

	c := cron.New()
	c.AddFunc("*/15 * * * *", func() {
		log.Println("Running scheduled scraper...")
		scraper.Run(q)
	})
	c.Start()

	select {}
}
