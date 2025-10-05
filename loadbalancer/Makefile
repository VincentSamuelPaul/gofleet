# Build binaries
build-loadbalancer:
	go build -o bin/loadbalancer ./loadbalancer

build-node1:
	go build -o bin/node1 ./nodes/node1

build-node2:
	go build -o bin/node2 ./nodes/node2

build-node3:
	go build -o bin/node3 ./nodes/node3

build-tests:
	cd tests && go build -o ../bin/loadtest

build: build-loadbalancer build-node1 build-node2 build-node3 build-tests

# Run services
run-loadbalancer:
	./bin/loadbalancer

run-node1:
	./bin/node1

run-node2:
	./bin/node2

run-node3:
	./bin/node3

run-tests:
	./bin/loadtest

# Start all services in parallel (foreground)
start-all: build
	@echo "Starting node1"
	@./bin/node1 &

	@echo "Starting node2"
	@./bin/node2 &

	@echo "Starting node3"
	@./bin/node3 &

	@echo "Starting load balancer"
	@./bin/loadbalancer &

	@echo "All services started."

# Stop all running services (by name)
stop-all:
	@echo "Stopping all services..."
	@pkill -f bin/loadbalancer || true
	@pkill -f bin/node1 || true
	@pkill -f bin/node2 || true
	@pkill -f bin/node3 || true
	@echo "All services stopped."

# Clean build artifacts
clean:
	rm -rf bin
