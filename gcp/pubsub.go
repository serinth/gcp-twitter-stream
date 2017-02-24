package gcp

import (
	"log"

	"cloud.google.com/go/pubsub"
	configuration "github.com/serinth/gcp-twitter-stream/config"
	"golang.org/x/net/context"
)

func handler() {
	ctx := context.Background()

	config, _ := configuration.GetConfig()

	// TODO do something with client instead of _
	_, err := pubsub.NewClient(ctx, config.GCP.Project)

	if err != nil {
		log.Fatal("PubSub client could not be initialized")
	}

}
