#!/usr/bin/env bash
set -e

go build -mod=readonly -o bin/protoc-gen-go google.golang.org/protobuf/cmd/protoc-gen-go
