#!/bin/bash

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative projects.proto
protoc --go_out=../../queue-worker/projects --go_opt=paths=source_relative --go-grpc_out=../../queue-worker/projects --go-grpc_opt=paths=source_relative projects.proto
