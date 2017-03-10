#!/bin/bash
# gcloud beta emulators pubsub start

if [[ $# -eq 0 ]] ; then
    echo 'No pubsub emulator port specified'
    exit 1
fi

export PUBSUB_EMULATOR_HOST=localhost:$1
go build .
./subscriber.exe
