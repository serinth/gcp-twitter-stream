package gcp

import (
	"log"
	"strconv"

	"time"

	"cloud.google.com/go/pubsub"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/golang/protobuf/proto"
	configuration "github.com/serinth/gcp-twitter-stream/config"
	tweetpb "github.com/serinth/gcp-twitter-stream/protobuf"
	"golang.org/x/net/context"
)

// Publisher with Context and GCP PubSub client
type Publisher struct {
	Client  pubsub.Client
	context context.Context
	config  configuration.Config
}

// NewPublisher returns new PubSub
func NewPublisher(configPath string) *Publisher {
	ctx := context.Background()
	config, _ := configuration.GetConfig(configPath)
	pubsubClient, _ := pubsub.NewClient(ctx, config.GCP.Project)
	return &Publisher{Client: *pubsubClient, context: ctx, config: config}
}

// Send tweet to PubSub stream
func (p *Publisher) Send(tweet *twitter.Tweet) {
	log.Println("Publishing tweet message: ", tweet.Text)
	topic, err := ensureTopicExists(p, p.config.PubSub.Topic)

	if err != nil {
		log.Fatal("Failed to ensure topic is ready for streaming: ", err)
	}

	msg := &tweetpb.Tweet{
		TweetId:       strconv.FormatInt(tweet.ID, 10),
		IngestionDate: time.Now().Format("2006-01-02 15:04:05"),
		Name:          tweet.User.Name,
		Tweet:         tweet.Text,
	}

	serializedMessage, err := proto.Marshal(msg)

	if err != nil {
		log.Println("Error converting to binary for pub/sub: ", err)
	} else {
		_ = topic.Publish(p.context, &pubsub.Message{Data: serializedMessage})
	}

}

// EnsureTopicExists creates a PubSub Topic if one doesn't exist
var ensureTopicExists = func(p *Publisher, t string) (*pubsub.Topic, error) {
	var err error

	topic := p.Client.Topic(t)

	topicExists, _ := topic.Exists(p.context)

	if topicExists == false {
		topic, err = p.Client.CreateTopic(p.context, t)
		if err != nil {
			log.Println("Could not create Topic: ", t)
		}
	}
	return topic, err
}
