[![Build Status](https://travis-ci.org/serinth/gcp-twitter-stream.svg?branch=master)](https://travis-ci.org/serinth/gcp-twitter-stream)
# Overview

Example project to implement the following architecture on GCP:
```
                   +------------------------------+
                   |Kubernetes Cluster            |
                   |                              |
                   |                              |
                   |      +---------------+       |        +---------+
+-------------+    |      |               |       |        |         |
|  Twitter    +---------> |   Publisher   +--------------> | Pub/Sub |
+-------------+    |      |               |       |        |         |
                   |      +---------------+       |        +----+----+
                   |                              |             |
                   |                              |             |
                   |      +---------------+       |             |
                   |      |               | <-------------------+
                   |      |  Subscriber   |       |
                   |      |               |       |
                   |      +------+--------+       |        +------------+
                   |             |                |        |            |
                   |             +-----------------------> |  BigQuery  |
                   +------------------------------+        |            |
                                                           +------------+
```

# Requirements

- Google Cloud SDK and CLI
  - **beta** component installed
  - **pubsub-emulator** component installed
  - **cloud-datastore-emulator** installed
  - **kubectl** component installed
- Golang 1.7+
- Docker Native
- Google Cloud account enabled
- Twitter account
- Twitter app credentials (https://apps.twitter.com)


Replace <PORT> with the pubsub emulator port number.

# Workshop Instructions

This workshop will cover the following aspects of GCP:
- Getting started with a GCP project layout
- How to Dockerize a Golang app and push it to Google Container Registry (GCR)
- How to deploy apps onto Kubernetes on Google Container Engine (GKE)
- Publishing and Subscribing to tweets via Google Pub/Sub queue
- Putting data into BigQuery and querying from it

## Getting GCP Ready

1. Log onto GCP and create a new Project
2. Set the following variables for the default CLI options:

```bash
# Ensure we're logged in and everything is ready to go. Skip init if you've already done it
gcloud init
gcloud auth application-default login


gcloud config list # See whats there already
gcloud config set project PROJECT_NAME

# gcloud compute regions list
gcloud config set compute/region asia-northeast1

# gcloud compute zones list 
gcloud config set compute/zone asia-northeast1-a
```

3. Set Kubernetes config

```bash
cd ~
mkdir .kube
echo "" > .kube/config

# Set KUBECONFIG in your shell's environment variable to point to this file
export KUBECONFIG=~/.kube/config

# Windows users set this in your KUBECONFIG environment variable
# %USERPROFILE%\.kube\config
```

4. Spin up a cluster while we do other things

```bash
gcloud container clusters create dius-cluster --zone asia-northeast1-a --num-nodes 2
# When that completes run:
gcloud container clusters get-credentials dius-cluster
```

## Get The App Running Locally

1. Setup Twitter environment variables in `~/.profile` or equivalent:

```bash
# Twitter Creds
export TWITTER_CONSUMER_KEY= ...
export TWITTER_CONSUMER_SECRET= ...
export TWITTER_ACCESS_TOKEN= ...
export TWITTER_ACCESS_SECRET= ...
```

Also put these values into the `publisher/Dockerfile`.

2. Build locally:

```bash
go get ./...
cd publisher && go build .
# Should get publisher executable or publisher.exe for Windows
cd ../subscriber && go build .
# Should get subscriber executable or subscriber.exe for Windows
```

3. Start the pub/sub emulator

```bash
gcloud beta emulators pubsub start
```

This should give you a port # in which pubsub is running

4. Run the publisher to make sure we're getting tweets

```bash
cd publisher
./start <PORTNUMBER>
```
* Note: you may need to modify `start.sh` to reference the executable name

## Build the Docker Containers and Push to GCR

1. Statically link and build the executable for Linux

```bash
cd publisher
CGO_ENABLED=0 GOOS=linux go build -a --ldflags="-s" --installsuffix cgo -o publisher
```

2. Build and tag the Docker images

```bash
docker build -t trumplisher:v1 .
docker tag trumplisher:v1 asia.gcr.io/PROJECT_ID/trumplisher:v1
```

3. Repeat steps 1 and 2 for Subscriber but rename the docker images and use **subscriber** as the Go binary name.

4. Push the images to GCR

```bash
gcloud docker -- push asia.gcr.io/PROJECT_ID/trumplisher:v1
```

You can view it on the web console.

## Create the BigQuery Table

1. First create a dataset and add it to our bqueryrc:

```bash
bq mk trump_data
echo "dataset_id=trump_data" >> ~/.bigqueryrc
```

2. Now create the table:

```bash
bq mk --schema bq_tweets_schema.json -t trump_data.tweets
```

```bash
bq head -n 10 PROJECT_ID:trump_data.tweets
```

