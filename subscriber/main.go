package main

import (
	"log"

	configuration "github.com/serinth/gcp-twitter-stream/config"
	"github.com/serinth/gcp-twitter-stream/gcp"
)

func main() {

	appConfig, err := configuration.GetConfig("../config/config.json")
	if err != nil {
		log.Fatal("Could not get config from config.json")
	}

	subscriber := gcp.NewSubscriber("../config/config.json")

	subscriber.Subscribe(appConfig.PubSub.Topic)
	subscriber.ListenAndHandle()

}
