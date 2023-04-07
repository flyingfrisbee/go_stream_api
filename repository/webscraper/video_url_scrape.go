package webscraper

import (
	"fmt"
	env "go_stream_api/environment"
	"regexp"

	"github.com/gocolly/colly"
)

func ScrapeVideoURL(episodeEndpoint string) string {
	var videoURL string

	c := colly.NewCollector()

	c.OnHTML(env.VideoURLSelector, func(e *colly.HTMLElement) {
		videoURL = e.Attr("src")
		regex := regexp.MustCompile(`[a-zA-Z0-9]`)
		startIndex := regex.FindStringIndex(videoURL)[0]

		// Invalid URL
		if startIndex != 0 {
			videoURL = fmt.Sprintf(`https://%s`, videoURL[startIndex:])
		}
	})

	c.OnError(errorCallback)

	url := env.BaseURLForScraping + episodeEndpoint
	c.Visit(url)

	return videoURL
}
