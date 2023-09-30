#!/bin/bash

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative queue.proto
protoc --go_out=../../queue-listener/queue --go_opt=paths=source_relative --go-grpc_out=../../queue-listener/queue --go-grpc_opt=paths=source_relative queue.proto