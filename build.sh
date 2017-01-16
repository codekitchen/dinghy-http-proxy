#!/bin/bash
set -e

docker run --rm -v "$PWD":/go/src/github.com/codekitchen/dinghy-http-proxy -w /go/src/github.com/codekitchen/dinghy-http-proxy golang:1.8 go build -v -o join-networks
tar czvf join-networks.tar.gz join-networks
rm join-networks
