package gcp

import (
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/dghubble/go-twitter/twitter"
	configuration "github.com/serinth/gcp-twitter-stream/config"
	"golang.org/x/net/context"
)

var config, _ = configuration.GetConfig()

// PubSub with Context and GCP PubSub client
type PubSub struct {
	Client  pubsub.Client
	context context.Context
}

// NewPubSub returns new PubSub
func NewPubSub() *PubSub {
	ctx := context.Background()
	pubsubClient, _ := pubsub.NewClient(ctx, config.GCP.Project)
	return &PubSub{Client: *pubsubClient, context: ctx}
}

// Send tweet to PubSub stream
func (p *PubSub) Send(tweet *twitter.Tweet) {
	log.Println("Publishing tweet message: ", tweet.Text)
	topic, err := ensureTopicExists(p, config.PubSub.Topic)
	if err != nil {
		log.Fatal("Failed to ensure topic is ready for streaming.")
	}
	_, err = topic.Publish(p.context, &pubsub.Message{Data: []byte(tweet.Text)})
	if err != nil {
		log.Println("Failed to send to topic witih error: ", err)
	}
}

// EnsureTopicExists creates a PubSub Topic if one doesn't exist
var ensureTopicExists = func(p *PubSub, t string) (*pubsub.Topic, error) {
	var err error
	topic := p.Client.Topic(t)

	if exists, _ := topic.Exists(p.context); exists == false {
		topic, err = p.Client.CreateTopic(p.context, t)
		if err != nil {
			log.Println("Could not create Topic: ", t)
		}
	}
	return topic, err
}
