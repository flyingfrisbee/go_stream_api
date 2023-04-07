package main

import (
	"go_stream_api/api"
	env "go_stream_api/environment"
)

func main() {
	env.LoadEnvVariables()
	api.Run()
}
