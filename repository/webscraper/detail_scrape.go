package webscraper

import (
	"fmt"
	env "go_stream_api/environment"
	"go_stream_api/repository/database/domain"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func scrapeDetail(a *domain.Anime, url string) {
	c := colly.NewCollector()

	c.OnHTML(env.DetailSelector, func(e *colly.HTMLElement) {
		switch {
		case strings.Contains(e.Text, env.Keyword1):
			movieType := strings.TrimSpace(strings.Replace(e.Text, env.Keyword1, "", 1))
			a.Type = movieType
		case strings.Contains(e.Text, env.Keyword2):
			summary := strings.TrimSpace(strings.Replace(e.Text, env.Keyword2, "", 1))
			a.Summary = summary
		case strings.Contains(e.Text, env.Keyword3):
			genre := strings.TrimSpace(strings.Replace(e.Text, env.Keyword3, "", 1))
			a.Genre = genre
		case strings.Contains(e.Text, env.Keyword4):
			airingYear := strings.TrimSpace(strings.Replace(e.Text, env.Keyword4, "", 1))
			a.AiringYear = airingYear
		case strings.Contains(e.Text, env.Keyword5):
			status := strings.TrimSpace(strings.Replace(e.Text, env.Keyword5, "", 1))
			a.Status = status
		default:
		}
	})

	c.OnError(errorCallback)

	c.Visit(url)

	if a.ID == 0 {
		c = colly.NewCollector()
		c.OnHTML(env.IDAtDetailSelector, func(e *colly.HTMLElement) {
			idString := e.Attr(env.SuperSecretKey2)
			if idString != "" {
				id, err := strconv.Atoi(idString)
				if err != nil {
					log.Fatal(err)
				}

				a.ID = id
			}
		})
		c.OnError(errorCallback)
		c.Visit(url)

		c = colly.NewCollector()
		c.OnHTML(env.ImageURLAtDetailSelector, func(e *colly.HTMLElement) {
			imageURL := e.Attr("src")
			a.ImageURL = imageURL
		})
		c.OnError(errorCallback)
		c.Visit(url)
	}
}

func scrapeEpisodes(a *domain.Anime, url string) {
	c := colly.NewCollector()

	c.OnHTML(env.EpisodesSelector, func(e *colly.HTMLElement) {
		episodeText := strings.Replace(e.ChildText("div:first-child"), "EP ", "", 1)
		endpoint := strings.TrimSpace(e.Attr("href"))
		episode := domain.Episode{
			Text:     episodeText,
			Endpoint: endpoint,
		}

		a.Episodes = append(a.Episodes, episode)
	})

	c.OnError(errorCallback)

	c.Visit(url)
}

func scrapeStream(a *domain.Anime, url string) {
	c := colly.NewCollector()

	c.OnHTML(env.StreamSelector, func(e *colly.HTMLElement) {
		idString := e.ChildAttr(env.SuperSecretKey1, env.SuperSecretKey2)
		detailEndpoint := e.ChildAttr(".anime_video_body_cate > .anime-info > a", "href")

		if idString != "" {
			id, err := strconv.Atoi(idString)
			if err != nil {
				log.Fatal(err)
			}

			a.ID = id
		}

		if detailEndpoint != "" {
			a.DetailEndpoint = detailEndpoint
		}
	})

	c.OnError(errorCallback)

	c.Visit(url)
}

func reverseEpisodesOrder(a *domain.Anime) {
	result := []domain.Episode{}

	length := len(a.Episodes)
	for i := (length - 1); i >= 0; i-- {
		result = append(result, a.Episodes[i])
	}

	a.Episodes = result
}

func calculateUpdateTime(a *domain.Anime, baseTime time.Time) {
	currentTime := time.Now().UTC()
	timeDiff := currentTime.Sub(baseTime)

	resultTime := baseTime.Add(-timeDiff)
	a.UpdatedAt = resultTime
}

// If user click title from search bar result, use this scrape func
func ScrapeDetailAlternative(title, endpoint string) domain.Anime {
	anime := domain.Anime{Title: title}
	url := env.BaseURLForScraping + endpoint
	scrapeDetail(&anime, url)

	url = fmt.Sprintf(env.EpisodesURLFormat, anime.ID)
	scrapeEpisodes(&anime, url)

	// Necessary because the order from scraping is descending.
	// Ascending is preferable hence the function call
	reverseEpisodesOrder(&anime)

	return anime
}
