#!/bin/bash

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs.proto
protoc --go_out=../../broker-service/logs --go_opt=paths=source_relative --go-grpc_out=../../broker-service/logs --go-grpc_opt=paths=source_relative logs.proto
protoc --go_out=../../queue-worker/logs --go_opt=paths=source_relative --go-grpc_out=../../queue-worker/logs --go-grpc_opt=paths=source_relative logs.proto
