package gcp

import (
	"log"

	"cloud.google.com/go/bigquery"
	tweetpb "github.com/serinth/gcp-twitter-stream/protobuf"
)

type tweetRow struct {
	tweetID       string
	name          string
	tweet         string
	ingestionDate string
}

// Must implement bigquery ValueSaver interface
func (t *tweetRow) Save() (row map[string]bigquery.Value, insertID string, err error) {
	return map[string]bigquery.Value{
		"tweet_id":       t.tweetID,
		"name":           t.name,
		"tweet":          t.tweet,
		"ingestion_date": t.ingestionDate,
	}, t.tweetID, nil
}

func (s *Subscriber) insertRow(tweet *tweetpb.Tweet) {
	u := s.BQClient.Dataset(s.config.BigQuery.DatasetID).Table(s.config.BigQuery.TableID).Uploader()

	rows := []*tweetRow{
		{
			tweetID:       tweet.TweetId,
			name:          tweet.Name,
			tweet:         tweet.Tweet,
			ingestionDate: tweet.IngestionDate,
		},
	}

	if err := u.Put(s.context, rows); err != nil {
		log.Println("Failed to insert record: ", err)
	}

}
