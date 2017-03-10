package gcp

import (
	"log"

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
}

// NewSubscriber creates a subscriber class to listen on a topic
func NewSubscriber(configPath string) *Subscriber {
	ctx := context.Background()
	config, _ := configuration.GetConfig(configPath)
	pubsubClient, _ := pubsub.NewClient(ctx, config.GCP.Project)
	return &Subscriber{Client: *pubsubClient, context: ctx, config: config}
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
	it, err := s.subcription.Pull(s.context)
	if err != nil {
		log.Fatal("Could not get message iterator from subscription")
	}

	defer it.Stop()

	for {
		message, err := it.Next()

		if err != nil {
			break
		}

		tweet := &tweetpb.Tweet{}

		err = proto.Unmarshal(message.Data, tweet)

		if err != nil {
			log.Println("Failed to deserialize message: ", err)
		} else {
			log.Println("Got Message: ", tweet.String())
		}

		message.Done(true)
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
