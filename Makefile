LOGGER_BINARY=loggerService
BROKER_BINARY=brokerService

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

build_logger:
	@echo "Building logger binary..."
	cd logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./api
	chmod +x ./logger-service/${LOGGER_BINARY}
	@echo "Done!"

build_broker:
	@echo "Building broker binary..."
	cd broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./api
	chmod +x ./broker-service/${BROKER_BINARY}
	@echo "Done!"