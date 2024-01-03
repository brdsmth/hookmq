ROOT_DIR := $(shell pwd)

api:    
	cd cmd/api && ROOT_DIR=$(ROOT_DIR) go run main.go &

runner:
	cd cmd/runner && ROOT_DIR=$(ROOT_DIR) go run main.go &

scheduler:
	cd cmd/scheduler && ROOT_DIR=$(ROOT_DIR) go run main.go &

start: api runner scheduler
	@echo "All services have started..."
