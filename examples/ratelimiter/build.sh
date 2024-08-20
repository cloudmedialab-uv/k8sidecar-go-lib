#!/bin/bash

TAG="1.0.0"
REPO_RATELIMITER=cloudmedialab/sidecar-ratelimiter

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo $SCRIPT_DIR

docker build $SCRIPT_DIR -t k8sidecar/examples/ratelimiter:$TAG -f $SCRIPT_DIR/Dockerfile

docker tag k8sidecar/examples/ratelimiter:$TAG $REPO_RATELIMITER:$TAG

docker push $REPO_RATELIMITER:$TAG
