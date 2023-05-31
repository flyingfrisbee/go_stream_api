package webscraper

import (
	"fmt"
	env "go_stream_api/environment"
	"time"

	"github.com/gocolly/colly"
)

type TitleSearchResult struct {
	Title    string `json:"title"`
	Endpoint string `json:"endpoint"`
}

func ScrapeAnimeTitlesByKeyword(keyword string) []TitleSearchResult {
	result := []TitleSearchResult{}

	c := createNewCollectorWithCustomTimeout(1 * time.Minute)

	c.OnHTML(env.TitlesSelector, func(e *colly.HTMLElement) {
		title := e.Text
		endpoint := e.Attr("href")
		searchResult := TitleSearchResult{
			Title:    title,
			Endpoint: endpoint,
		}

		result = append(result, searchResult)
	})

	c.OnError(errorCallback)

	url := fmt.Sprintf(env.TitleSearchURLFormat, env.BaseURLForScraping, keyword)
	c.Visit(url)

	return result
}
