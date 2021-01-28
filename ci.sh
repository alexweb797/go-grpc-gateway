#!/bin/bash

docker run --rm  \
  -v $(pwd):/go/src \
  -e MAKE_TARGET="$1" \
  -e GOOS=linux \
  -e GOARCH=amd64 \
  golang:1.14 \
  ./src/build.sh
