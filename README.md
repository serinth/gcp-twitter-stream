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
- Golang 1.7+
- Docker Native
- Google Cloud account enabled
- Twitter account
- Twitter app credentials (https://apps.twitter.com)

## Optional Install - Protocol Buffers 

https://developers.google.com/protocol-buffers/

```bash
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```


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
2. Ensure that `Enable Billing` has been clicked in Billing or Compute Engine
3. Set the following variables for the default CLI options:

```bash
# Ensure we're logged in and everything is ready to go. Skip init if you've already done it
gcloud init

# Install extra components for Kubernetes and Pub/Sub emulation
gcloud components install beta
gcloud components install pubsub-emulator
gcloud components install kubectl

gcloud auth application-default login

gcloud config list # See whats there already
gcloud config set project PROJECT_ID

# gcloud compute regions list
gcloud config set compute/region asia-northeast1

# gcloud compute zones list 
gcloud config set compute/zone asia-northeast1-a
```

4. Set Kubernetes config

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
gcloud container clusters create dius-cluster --zone asia-northeast1-a --num-nodes 2 --scopes=compute-rw,monitoring,logging-write,storage-rw,bigquery,https://www.googleapis.com/auth/pubsub
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

Also put these values into the `PublisherDockerfile`.

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
./start.sh <PORTNUMBER>
```
* Note: you may need to modify `start.sh` to reference the executable name

5. Lastly modify the `config/config.json` and put in your Project ID

## Build the Docker Containers and Push to GCR

1. Statically link and build the executable for Linux

```bash
cd publisher
CGO_ENABLED=0 GOOS=linux go build -a --ldflags="-s" --installsuffix cgo -o publisher
```

2. Build and tag the Docker images

```bash
cd ..
docker build -f PublisherDockerfile -t trumplisher:v1 .
docker tag trumplisher:v1 asia.gcr.io/PROJECT_ID/trumplisher:v1
```
3. Statically link the subscriber

```bash
cd subscriber
CGO_ENABLED=0 GOOS=linux go build -a --ldflags="-s" --installsuffix cgo -o subscriber
cd ..
docker build -f SubscriberDockerfile -t trumpscriber:v1 .
docker tag trumpscriber:v1 asia.gcr.io/PROJECT_ID/trumpscriber:v1
```


4. Push the images to GCR

```bash
gcloud docker -- push asia.gcr.io/PROJECT_ID/trumplisher:v1
gcloud docker -- push asia.gcr.io/PROJECT_ID/trumpscriber:v1
```

You can view it on the web console.

## Create the Pub/Sub Topic

1. Navigate to the console and **Enable the API**
2. Create the *trumpisms* topic

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

## Deploy Publisher and Subscriber to Kubernetes

1. Modify `publisher.yaml` and `subscriber.yaml` to include the project id in the container to run

2. Deploy the publisher

```bash
kubectl create -f publisher.yaml
```

3. Deploy the subscriber

```bash
kubectl create -f subscriber.yaml
```

## Query Data from CLI

```bash
bq head -n 10 --format=prettyjson PROJECT_ID:trump_data.tweets
```

## Write SQL Queries on BigQuery Console

1. Select the tweets table and start writing SQL Queries

# Horizontal Pod Autoscaling (more pods)

1. Delete existing resources

a. via CLI:

```bash
kubectl delete deployment trumplisher-deployment
kubectl delete deployment trumpscriber-deployment
gcloud container clusters delete dius-cluster
```

b. via console:
- Just go to delete the entire cluster

2. Spin up a new cluster with 1 node and max 2 nodes and enable autoscaling

```bash
gcloud alpha container clusters create dius-cluster --zone asia-northeast1-a --enable-autoscaling --num-nodes 1 --min-nodes 1 --max-nodes 2 --scopes=compute-rw,monitoring,logging-write,storage-rw,bigquery 

gcloud container clusters get-credentials dius-cluster
```

# Cluster Scaling (more nodes -- upscaling hardware and access rights)

1. Create a stress application with `StressDockerfile`

```bash
docker build -f StressDockerfile -t stress:v1 .
docker tag stress:v1 asia.gcr.io/PROJECT_ID/stress:v1
gcloud docker -- push asia.gcr.io/PROJECT_ID/stress:v1
```
2. Generate some load:

Edit the stress.yaml file to add in your project id then;

```bash
kubectrl create -f stress.yaml
```

3. Use `kubectl get nodes` and cloud console to see what happens.

# Horizontal Pod Autoscaling (more pods)

1. Apply horizontal pod autoscaler:

```bash
kubectl create -f autoscale.yaml
```

2. Use a combination of:

```bash
kubectl get pods
kubectl get deployments
kubectl get hpa
```

To see things scale.
