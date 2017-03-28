package gcp

import (
	"log"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"github.com/golang/protobuf/proto"
	configuration "github.com/serinth/gcp-twitter-stream/config"
	tweetpb "github.com/serinth/gcp-twitter-stream/protobuf"
	"golang.org/x/net/context"
)

// Subscriber is a GCP PubSub subscriber
type Subscriber struct {
	Client      pubsub.Client
	context     context.Context
	subcription pubsub.Subscription
	config      configuration.Config
	BQClient    bigquery.Client
}

// NewSubscriber creates a subscriber class to listen on a topic
func NewSubscriber(configPath string) *Subscriber {
	ctx := context.Background()
	config, _ := configuration.GetConfig(configPath)
	pubsubClient, err := pubsub.NewClient(ctx, config.GCP.Project)

	if err != nil {
		log.Fatal("Could not create pubsub Client with error: ", err)
	}

	bqClient, bqerr := bigquery.NewClient(ctx, config.GCP.Project)

	if bqerr != nil {
		log.Fatal("Could not create BigQuery Client with error: ", bqerr)
	}

	return &Subscriber{Client: *pubsubClient, context: ctx, config: config, BQClient: *bqClient}
}

// Subscribe to a topic. Creates a subscription if one doesn't exist
func (s *Subscriber) Subscribe(topic string) {

	pubsubTopic := s.getTopic(topic)
	if pubsubTopic == nil {
		log.Fatal("Topic does not exist")
	}

	var sub = s.getSubscription(s.config.PubSub.Topic)
	var err error

	if sub == nil {
		log.Println("Creating subscription: ", s.config.PubSub.Topic)
		sub, err = s.Client.CreateSubscription(s.context, s.config.PubSub.Topic, pubsubTopic, 0, nil)
	}

	if err != nil {
		log.Fatal("Could not create subscription with error: ", err)
	}

	s.subcription = *s.Client.Subscription(topic)
}

// ListenAndHandle will start handling the messages in subscription
func (s *Subscriber) ListenAndHandle() {
	err := s.subcription.Receive(s.context, func(context context.Context, message *pubsub.Message) {

		tweet := &tweetpb.Tweet{}

		erro := proto.Unmarshal(message.Data, tweet)

		if erro != nil {
			log.Println("Failed to deserialize message: ", erro)
		} else {

			log.Println("Got Message: ", tweet.String())
			s.insertRow(tweet)
		}

		message.Ack()
	})
	if err != nil {
		log.Fatal("Could not get message iterator from subscription")
	}
}

func (s *Subscriber) getTopic(topic string) *pubsub.Topic {
	pubsubTopic := s.Client.Topic(s.config.PubSub.Topic)
	topicExists, err := pubsubTopic.Exists(s.context)

	if err == nil && topicExists {
		return pubsubTopic
	}

	return nil
}

func (s *Subscriber) getSubscription(name string) *pubsub.Subscription {
	iter := s.Client.Subscriptions(s.context)

	var result *pubsub.Subscription
	for {
		sub, err := iter.Next()
		if err != nil {
			log.Println("Could not find subscription: ", name)
			break
		}

		if sub.ID() == name {
			result = sub
			break
		}
	}

	return result
}
