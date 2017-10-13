#!/bin/bash
#
# Compiles the requisite .proto files into Go sources
#
set -euo pipefail

# Build the imports line
[ -d "${GOPATH}/src/github.com/google/protobuf/src" ] || go get -d -v -u github.com/google/protobuf/src || :
[ -d "${GOPATH}/src/github.com/googleapis/googleapis/google/rpc" ] || $(\
    go get -d -v -u github.com/googleapis/googleapis/google/rpc; \
    go get -d -v -u google.golang.org/genproto/googleapis/rpc/status ) || :
imports=".:proto:${GOPATH}/src/github.com/google/protobuf/src:${GOPATH}/src/github.com/googleapis/googleapis/"

# Build the Go package prefix for imported sources
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PREFIX=${DIR#${GOPATH}/src/}
p4infoMapping="p4/config/p4info.proto=$PREFIX/proto/p4/config"

# Generate Go sources
protoc --go_out=. proto/p4/tmp/p4config.proto
protoc -I=$imports --go_out=. proto/p4/config/p4info.proto
protoc -I=$imports --go_out=plugins=grpc,M${p4infoMapping}:. proto/p4/p4runtime.proto