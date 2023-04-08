package webscraper

import (
	"context"
	"fmt"
	env "go_stream_api/environment"
	db "go_stream_api/repository/database"
	"log"
	"time"

	"github.com/gocolly/colly"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func StartScrapingService() {
	ctx, cancel = context.WithCancel(context.Background())

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping scraping service...")
			return
		default:
			runScrapeLoop()
		}
	}
}

// Scrape data and save to database
func runScrapeLoop() {
	baseTime := time.Now().UTC()

	for i := 1; i <= 10; i++ {
		url := fmt.Sprintf("%s?page=%d", env.BaseURLForScraping, i)
		animes := scrapeHome(url)

		for _, anime := range animes {
			url = env.BaseURLForScraping + anime.StreamEndpoint
			scrapeStream(&anime, url)

			url = fmt.Sprintf(env.EpisodesURLFormat, anime.ID)
			scrapeEpisodes(&anime, url)

			url = env.BaseURLForScraping + anime.DetailEndpoint
			scrapeDetail(&anime, url)

			// Necessary because the order from scraping is descending.
			// Ascending is preferable hence the function call
			reverseEpisodesOrder(&anime)
			calculateUpdateTime(&anime, baseTime)

			err := db.Conn.UpsertAnime(&anime)
			if err != nil {
				log.Fatal(err)
			}

			// err := anime.ProcessAnimeData()
			// if err != nil {
			// 	log.Fatal(err)
			// }
		}
	}
}

func errorCallback(r *colly.Response, err error) {
	log.Printf(
		"error when visiting %s\nerror message: %s\n",
		r.Request.URL.String(),
		err.Error(),
	)
}
