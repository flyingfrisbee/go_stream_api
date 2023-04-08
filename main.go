package main

import (
	"go_stream_api/api"
	env "go_stream_api/environment"
	db "go_stream_api/repository/database"
	"go_stream_api/repository/webscraper"
)

func main() {
	env.LoadEnvVariables()
	db.StartConnectionToDB()
	api.Run()
	go webscraper.StartScrapingService()
}
