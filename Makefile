.PHONY: help run build test clean docker deploy logs status

BINARY=plug-monitor

help:
	@echo "Available commands:"
	@echo ""
	@echo "  make run      - Run locally"
	@echo "  make build    - Build binary"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Remove binary"
	@echo "  make docker   - Build Docker image"
	@echo "  make deploy   - Rebuild and start container"
	@echo "  make logs     - Show container logs"
	@echo "  make status   - Show running containers"

run:
	go run ./cmd/server

build:
	go build -o $(BINARY) ./cmd/server

test:
	go test ./...

clean:
	rm -f $(BINARY)

docker:
	@if ! command -v docker >/dev/null; then \
		echo "Docker is not installed."; \
		exit 1; \
	fi
	docker compose build

deploy:
	@if ! command -v docker >/dev/null; then \
		echo "Docker is not installed."; \
		exit 1; \
	fi
	docker compose up -d --build

logs:
	@if ! command -v docker >/dev/null; then \
		echo "Docker is not installed."; \
		exit 1; \
	fi
	docker compose logs -f

status:
	@if ! command -v docker >/dev/null; then \
		echo "Docker is not installed."; \
		exit 1; \
	fi
	docker compose ps