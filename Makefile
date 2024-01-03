.PHONY: api runner scheduler

api:	
	cd cmd/api && go run main.go &

runner:
	cd cmd/runner && go run main.go &

scheduler:
	cd cmd/scheduler && go run main.go &

start: api runner scheduler
	@echo "All services have started..."
