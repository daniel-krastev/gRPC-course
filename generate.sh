#!/bin/bash

protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. ./calc/calc_proto/calc.proto

protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. ./greet/greet_pb/greet.proto

# protoc -I ./calc/calc_proto/ ./calc/calc_proto/calc.proto --go_out=plugins=grpc:./calc/calc_proto/.
# protoc -I ./greet/greet_pb/ ./greet/greet_pb/greet.proto --go_out=plugins=grpc:./greet/greet_pb/.