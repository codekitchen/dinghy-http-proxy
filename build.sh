#!/bin/bash
set -e

docker run --rm -e GOOS=linux -e GOARCH=amd64 -e CGO_ENABLED=0 -v "$PWD":/go/src/github.com/codekitchen/dinghy-http-proxy -w /go/src/github.com/codekitchen/dinghy-http-proxy golang go get github.com/fsouza/go-dockerclient && env GOOS=linux GOARCH=amd64 go build -v -o join-networks
