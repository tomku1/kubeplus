#!/bin/bash

if (( $# < 1 )); then
    echo "./build-artifact.sh <latest | versioned>"
fi

artifacttype=$1

if [ "$artifacttype" = "latest" ]; then
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o crd-hook
    docker build --no-cache -t lmecld/pac-mutating-admission-webhook:latest .
    docker push lmecld/pac-mutating-admission-webhook:latest
fi

if [ "$artifacttype" = "versioned" ]; then
    version=`tail -1 versions.txt`
    echo "Building version $version"
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o crd-hook
    #docker build --no-cache -t lmecld/pac-mutating-admission-webhook:$version .
    #docker push lmecld/pac-mutating-admission-webhook:$version
    docker build --no-cache -t gcr.io/disco-horizon-103614/pac-mutating-admission-webhook:$version .
    docker push gcr.io/disco-horizon-103614/pac-mutating-admission-webhook:$version
fi



