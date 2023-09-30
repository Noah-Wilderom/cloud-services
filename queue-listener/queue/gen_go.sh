#!/bin/bash

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative queue.proto
protoc --go_out=../../queue-worker/queue --go_opt=paths=source_relative --go-grpc_out=../../queue-worker/queue --go-grpc_opt=paths=source_relative queue.proto
