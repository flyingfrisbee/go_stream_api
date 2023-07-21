package main

import (
	"go_stream_api/api"
	env "go_stream_api/environment"
	db "go_stream_api/repository/database"
	"go_stream_api/repository/webscraper"
)

var blockerCh chan struct{}

func main() {
	env.LoadEnvVariables()
	db.StartConnectionToDB()
	// webhook.StartWebhookService()
	go webscraper.StartScrapingService()
	api.Run()
	<-blockerCh
}

func init() {
	blockerCh = make(chan struct{})
}
