package webscraper

import (
	env "go_stream_api/environment"
	"strings"

	"github.com/gocolly/colly"
)

func scrapeHome(url string) []anime {
	animes := []anime{}

	c := colly.NewCollector()

	c.OnHTML(env.HomeSelector, func(e *colly.HTMLElement) {
		title := e.ChildText(".name > a")
		imgURL := e.ChildAttr("div > a > img", "src")
		latestEpisode := strings.Replace(e.ChildText(".episode"), "Episode ", "", 1)
		streamEndpoint := e.ChildAttr(".name > a", "href")

		anime := anime{
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
