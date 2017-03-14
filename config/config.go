package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config for application
type Config struct {
	PubSub   pubSubConfig   `json:"pubsub"`
	Twitter  twitterConfig  `json:"twitter"`
	GCP      gcpConfig      `json:"gcp"`
	BigQuery bigQueryConfig `json:"bigQuery"`
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

type bigQueryConfig struct {
	DatasetID string `json:"datasetId"`
	TableID   string `json:"tableId"`
}

// GetConfig returns system config
func GetConfig(filePath string) (Config, error) {
	raw, err := ioutil.ReadFile(filePath)

	var config Config
	err = json.Unmarshal(raw, &config)

	return config, err
}
