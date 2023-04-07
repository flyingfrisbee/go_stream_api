package main

import (
	"go_stream_api/api"
	env "go_stream_api/environment"
	"go_stream_api/repository/webscraper"
)

func main() {
	env.LoadEnvVariables()
	api.Run()
	go webscraper.StartScrapingService()
}
