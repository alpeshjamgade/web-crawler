URL_FRONTIER_BINARY=urlFrontierApp

## up: starts all containers in the background without forcing _build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_url_frontier
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

## build_url_frontier: builds the url_frontier binary as a linux executable
build_url_frontier:
	@echo "Building url frontier binary..."
	cd ./url-frontier && env GOOS=linux CGO_ENABLED=0 go build -o ./_build/${URL_FRONTIER_BINARY} ./cmd
	@echo "Done!"
