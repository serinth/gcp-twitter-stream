[![Build Status](https://travis-ci.org/serinth/gcp-twitter-stream.svg?branch=master)](https://travis-ci.org/serinth/gcp-twitter-stream)
# Overview

Example project to implement the following architecture on GCP:
```
             +------------------------+
             |                        |
             |  Trump Twitter Stream  |
             |           App          |
             |                        |
             +------------+-----------+
                          |
                          |
                  +-------v-------+
                  |               |
                  |  GCP Pub/Sub  |
                  |               |
                  +-------+-------+
                          |
                          |
                          |
               +----------v---------+
               |                    |
          +----+    GCP Dataflow    +----+
          |    |                    |    |
          |    +--------------------+    |
          |                              |
          |                              |
          v                              v
+---------+-------+               +------+------+
|                 |               |             |
|  Cloud Storage  | +---------->  |  Big Query  |
|                 |               |             |
+-----------------+               +-------------+
```

# Requirements

- Google Cloud SDK and CLI
  - **beta** component installed
  - **pubsub-emulator** component installed
  - **cloud-datastore-emulator** installed
- Golang 1.7+

# Quick Start

```bash
go get ./...
go build .
./start.sh <PORT>
```

Replace <PORT> with the pubsub emulator port number.
