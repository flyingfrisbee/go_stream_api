package webscraper

import (
	"context"
	"fmt"
	env "go_stream_api/environment"
	"go_stream_api/notification"
	db "go_stream_api/repository/database"
	"go_stream_api/repository/database/domain"
	"log"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
)

func StartScrapingService() {
	wg.Add(1)
	ctx, cancel = context.WithCancel(context.Background())

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping scraping service...")
			wg.Done()
			return
		default:
			runScrapeLoop()
		}
	}
}

// Will block until webscraper finish on going scraping process
func Stop() {
	cancel()
	wg.Wait()
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

			err := processAnime(&anime)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func processAnime(a *domain.Anime) error {
	// Scraper failed to get anime id or episodes, return early
	if a.ID == 0 || len(a.Episodes) == 0 {
		return nil
	}

	// Get enum by comparing latest episode
	latestEpisode, err := db.Conn.Pg.GetLatestEpisode(a.ID)
	if err != nil {
		return err
	}

	comparisonResult := domain.NewEpisodeFound
	if latestEpisode == nil {
		comparisonResult = domain.EntryNotFound
	} else if *latestEpisode == a.LatestEpisode {
		comparisonResult = domain.NoChangesFound
	}

	storageType := db.Postgres
	if len(a.Episodes) > 30 {
		storageType = db.MongoDB
	}

	if storageType == db.Postgres {
		// Postgres
		switch comparisonResult {
		case domain.EntryNotFound:
			err = db.Conn.Pg.UpsertAnime(a)
			if err != nil {
				return err
			}
			rowsInserted, err := db.Conn.Pg.InsertEpisodesToPostgres(a, 0)
			if err != nil {
				return err
			}
			if rowsInserted < 1 {
				return fmt.Errorf("failed inserting episodes of anime %s", a.Title)
			}
		case domain.NewEpisodeFound:
			err = db.Conn.Pg.UpsertAnime(a)
			if err != nil {
				return err
			}
			epsCount, err := db.Conn.Pg.GetEpisodesCount(a)
			if err != nil {
				return err
			}
			newEpsCount := len(a.Episodes) - epsCount
			if newEpsCount >= 1 {
				rowsInserted, err := db.Conn.Pg.InsertEpisodesToPostgres(a, newEpsCount)
				if err != nil {
					return err
				}
				if rowsInserted < 1 {
					return fmt.Errorf("failed inserting episodes of anime %s", a.Title)
				}
			}
		case domain.NoChangesFound:
		}
	} else {
		// MongoDB
		switch comparisonResult {
		case domain.EntryNotFound:
			err = db.Conn.Pg.UpsertAnime(a)
			if err != nil {
				return err
			}
			episodes := a.GetEpisodesAsSliceInterface()
			err = db.Conn.Mongo.InsertEpisodes(a, episodes, 0)
			if err != nil {
				return err
			}
		case domain.NewEpisodeFound:
			err = db.Conn.Pg.UpsertAnime(a)
			if err != nil {
				return err
			}
			episodes := a.GetEpisodesAsSliceInterface()
			epsCount, err := db.Conn.Mongo.GetEpisodesCount(a.ID)
			if err != nil {
				return err
			}
			newEpsCount := len(episodes) - epsCount
			if newEpsCount >= 1 {
				err = db.Conn.Mongo.InsertEpisodes(a, episodes, newEpsCount)
				if err != nil {
					return err
				}
			}
		case domain.NoChangesFound:
		}
	}

	err = notification.SendNotificationMessageToUsers(comparisonResult, a)
	if err != nil {
		return err
	}

	return nil
}

func errorCallback(r *colly.Response, err error) {
	log.Printf(
		"error when visiting %s\nerror message: %s\n",
		r.Request.URL.String(),
		err.Error(),
	)
}
