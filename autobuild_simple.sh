#!/usr/bin/env bash

CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build main.go
scp main root@yourhost:/root/im/im_go_cent