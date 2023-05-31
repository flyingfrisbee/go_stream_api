package webscraper

import (
	env "go_stream_api/environment"
	"go_stream_api/repository/database/domain"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func scrapeHome(url string) []domain.Anime {
	animes := []domain.Anime{}

	c := createNewCollectorWithCustomTimeout(1 * time.Minute)

	c.OnHTML(env.HomeSelector, func(e *colly.HTMLElement) {
		title := e.ChildText(".name > a")
		imgURL := e.ChildAttr("div > a > img", "src")
		latestEpisode := strings.Replace(e.ChildText(".episode"), "Episode ", "", 1)
		streamEndpoint := e.ChildAttr(".name > a", "href")

		anime := domain.Anime{
			Title:          title,
			ImageURL:       imgURL,
			LatestEpisode:  latestEpisode,
			StreamEndpoint: streamEndpoint,
		}

		animes = append(animes, anime)
	})

	c.OnError(errorCallback)

	c.Visit(url)

	return animes
}
