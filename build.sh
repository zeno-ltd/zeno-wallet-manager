#!/bin/bash


CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -installsuffix cgo -o kms .

docker build -t zeno/zeno-wallet-manager -f Dockerfile.minimal .
