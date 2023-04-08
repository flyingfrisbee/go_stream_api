package main

import (
	"go_stream_api/api"
	env "go_stream_api/environment"
	db "go_stream_api/repository/database"
	"go_stream_api/repository/webscraper"
	"go_stream_api/webhook"
)

var blockerCh chan struct{}

func main() {
	env.LoadEnvVariables()
	db.StartConnectionToDB()
	webhook.StartWebhookService()
	api.Run()
	<-blockerCh
	go webscraper.StartScrapingService()
}

func init() {
	blockerCh = make(chan struct{})
}
