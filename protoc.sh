#!/bin/bash

protoc -I=./api/proto --go_out=plugins=grpc:./pkg/model ./api/proto/**/*.proto
