package webscraper

import (
	"context"
	"fmt"
	env "go_stream_api/environment"
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
			anime.scrapeStream(url)

			url = fmt.Sprintf(env.EpisodesURLFormat, anime.ID)
			anime.scrapeEpisodes(url)

			url = env.BaseURLForScraping + anime.DetailEndpoint
			anime.scrapeDetail(url)

			// Necessary because the order from scraping is descending.
			// Ascending is preferable hence the function call
			anime.reverseEpisodesOrder()
			anime.calculateUpdateTime(baseTime)

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
