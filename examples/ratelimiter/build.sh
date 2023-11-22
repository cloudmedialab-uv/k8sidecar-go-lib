#!/bin/bash

TAG="1.0.0.test"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

docker build $SCRIPT_DIR -t sidecar/filter/controller:$TAG -f $SCRIPT_DIR/Dockerfile
