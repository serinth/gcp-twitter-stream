package gcp

import (
	"log"

	"cloud.google.com/go/pubsub"
	configuration "github.com/serinth/gcp-twitter-stream/config"
	"golang.org/x/net/context"
)

func Handler() {
	ctx := context.Background()

	config, _ := configuration.GetConfig()

	client, err := pubsub.NewClient(ctx, config.GCP.Project)

	if err != nil {
		log.Fatal("PubSub client could not be initialized")
	}

}
