package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config for application
type Config struct {
	PubSub  pubSubConfig  `json:"pubsub"`
	Twitter twitterConfig `json:"twitter"`
	GCP     gcpConfig     `json:"gcp"`
}

type pubSubConfig struct {
	Topic string `json:"topic"`
}

type twitterConfig struct {
	Track []string `json:"track"`
}

type gcpConfig struct {
	Project string `json:"project"`
}

// GetConfig returns system config
func GetConfig() (Config, error) {
	cwd, err := os.Getwd()
	raw, err := ioutil.ReadFile(cwd + "/config/config.json")

	var config Config
	err = json.Unmarshal(raw, &config)

	return config, err
}
