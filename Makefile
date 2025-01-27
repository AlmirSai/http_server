# Variables
DOCKER_COMPOSE = docker-compose
GOLANG_CI = golangci-lint
GO = go
KUBECTL = kubectl
COVER_PROFILE = coverage.out
BENCH_PROFILE = bench.out
DOCKER_REGISTRY = your-registry
APP_VERSION = 0.1.0

# Docker commands
.PHONY: up
up:
	$(DOCKER_COMPOSE) up -d

.PHONY: down
down:
	$(DOCKER_COMPOSE) down

.PHONY: build
build:
	$(DOCKER_COMPOSE) build

.PHONY: logs
logs:
	$(DOCKER_COMPOSE) logs -f

.PHONY: ps
ps:
	$(DOCKER_COMPOSE) ps

# Development commands
.PHONY: test
test:
	$(GO) test ./... -v -race -coverprofile=$(COVER_PROFILE)
	go tool cover -html=$(COVER_PROFILE)

.PHONY: test-short
test-short:
	$(GO) test ./... -v -short

.PHONY: bench
bench:
	$(GO) test -bench=. -benchmem ./... > $(BENCH_PROFILE)

.PHONY: test-integration
test-integration:
	$(GO) test ./... -v -tags=integration

.PHONY: coverage-report
coverage-report:
	$(GO) test ./... -coverprofile=$(COVER_PROFILE)
	go tool cover -func=$(COVER_PROFILE)

.PHONY: lint
lint:
	$(GOLANG_CI) run

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: vendor
vendor:
	$(GO) mod vendor

# Kubernetes commands
.PHONY: k8s-apply
k8s-apply:
	$(KUBECTL) apply -f k8s/base/namespace.yaml
	$(KUBECTL) apply -f k8s/base/config-map.yaml
	$(KUBECTL) apply -f k8s/base/postgres-secret.yaml
	$(KUBECTL) apply -f k8s/base/logging.yaml

.PHONY: k8s-delete
k8s-delete:
	$(KUBECTL) delete -f k8s/base

# SSL Certificate generation
.PHONY: ssl-cert
ssl-cert:
	mkdir -p certs
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
		-keyout certs/server.key -out certs/server.crt

# Clean commands
.PHONY: clean
clean:
	$(DOCKER_COMPOSE) down -v
	rm -rf bin/*
	rm -rf vendor/
	rm -f $(COVER_PROFILE)
	rm -f $(BENCH_PROFILE)

.PHONY: clean-images
clean-images:
	docker rmi $$(docker images -q) -f

.PHONY: clean-volumes
clean-volumes:
	docker volume rm $$(docker volume ls -q)

.PHONY: clean-all
clean-all: clean
	rm -rf certs/*
	docker system prune -af

# Migration commands
.PHONY: migrate-up
migrate-up:
	go run cmd/migrate/main.go up

.PHONY: migrate-down
migrate-down:
	go run cmd/migrate/main.go down

# Service specific commands
.PHONY: run-auth
run-auth:
	go run ./auth-service/cmd/app/main.go

.PHONY: run-user
run-user:
	go run ./user-service/cmd/app/main.go

.PHONY: run-post
run-post:
	go run ./post-service/cmd/app/main.go

.PHONY: run-media
run-media:
	go run ./media-service/cmd/app/main.go

.PHONY: run-gateway
run-gateway:
	go run ./api-gateway/cmd/app/main.go

# PostgreSQL and Auth Service commands
.PHONY: postgres-auth
postgres-auth:
	docker run --name postgres-auth \
		-e POSTGRES_USER=auth_user \
		-e POSTGRES_PASSWORD=auth_password \
		-e POSTGRES_DB=auth_db \
		-p 5432:5432 \
		-d postgres:latest

.PHONY: run-auth-local
run-auth-local: postgres-auth
	export DB_HOST=localhost && \
	export DB_PORT=5432 && \
	export DB_USER=auth_user && \
	export DB_PASSWORD=auth_password && \
	export DB_NAME=auth_db && \
	export JWT_SECRET_KEY=your_jwt_secret_key && \
	go run ./auth-service/cmd/app/main.go

.PHONY: stop-postgres-auth
stop-postgres-auth:
	docker stop postgres-auth
	docker rm postgres-auth

.PHONY: run-gateway
run-gateway:
	go run ./api-gateway/cmd/app/main.go

# Kafka commands
.PHONY: kafka-up
kafka-up:
	$(DOCKER_COMPOSE) up -d zookeeper kafka-1 kafka-2

.PHONY: kafka-down
kafka-down:
	$(DOCKER_COMPOSE) stop zookeeper kafka-1 kafka-2
	$(DOCKER_COMPOSE) rm -f zookeeper kafka-1 kafka-2

.PHONY: kafka-topics
kafka-topics:
	docker exec kafka-1 kafka-topics --bootstrap-server localhost:9092 --list

.PHONY: kafka-create-topic
kafka-create-topic:
	@read -p "Enter topic name: " topic_name; \
	docker exec kafka-1 kafka-topics --bootstrap-server localhost:9092 --create --topic $$topic_name --partitions 3 --replication-factor 2

# MinIO commands
.PHONY: minio-up
minio-up:
	$(DOCKER_COMPOSE) up -d minio

.PHONY: minio-down
minio-down:
	$(DOCKER_COMPOSE) stop minio
	$(DOCKER_COMPOSE) rm -f minio

.PHONY: minio-mc
minio-mc:
	docker run -it --network backend-network --entrypoint /bin/sh minio/mc

# Database operations
.PHONY: db-backup
db-backup:
	docker exec postgres pg_dump -U $(DB_USER) $(DB_NAME) > backup.sql

.PHONY: db-restore
db-restore:
	docker exec -i postgres psql -U $(DB_USER) $(DB_NAME) < backup.sql

# Monitoring commands
.PHONY: monitor-services
monitor-services:
	docker stats

.PHONY: check-health
check-health:
	@echo "Checking services health..."
	curl -f http://localhost:8443/health || echo "API Gateway: DOWN"
	curl -f http://localhost:8080/health || echo "Auth Service: DOWN"

# Docker registry commands
.PHONY: docker-tag
docker-tag:
	docker tag api-gateway:latest $(DOCKER_REGISTRY)/api-gateway:$(APP_VERSION)
	docker tag auth-service:latest $(DOCKER_REGISTRY)/auth-service:$(APP_VERSION)

.PHONY: docker-push
docker-push: docker-tag
	docker push $(DOCKER_REGISTRY)/api-gateway:$(APP_VERSION)
	docker push $(DOCKER_REGISTRY)/auth-service:$(APP_VERSION)

# Environment setup
.PHONY: dev-env
dev-env:
	cp .env.dev .env

.PHONY: prod-env
prod-env:
	cp .env.prod .env

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "Development:"
	@echo "  make test            - Run all tests with coverage"
	@echo "  make test-short      - Run short tests"
	@echo "  make bench           - Run benchmarks"
	@echo "  make test-integration - Run integration tests"
	@echo "  make coverage-report - Generate coverage report"
	@echo "\nDatabase:"
	@echo "  make db-backup       - Backup PostgreSQL database"
	@echo "  make db-restore      - Restore PostgreSQL database"
	@echo "\nMonitoring:"
	@echo "  make monitor-services - Monitor Docker services"
	@echo "  make check-health    - Check services health"
	@echo "\nDocker:"
	@echo "  make docker-tag      - Tag Docker images"
	@echo "  make docker-push     - Push Docker images to registry"
	@echo "\nEnvironment:"
	@echo "  make dev-env         - Set up development environment"
	@echo "  make prod-env        - Set up production environment"
	@echo "  make up              - Start all services"
	@echo "  make down            - Stop all services"
	@echo "  make build           - Build all services"
	@echo "  make test            - Run tests"
	@echo "  make lint            - Run linter"
	@echo "  make clean           - Clean up build artifacts"
	@echo "  make k8s-apply      - Apply Kubernetes configurations"
	@echo "  make k8s-delete     - Delete Kubernetes configurations"
	@echo "  make ssl-cert       - Generate SSL certificates"
	@echo "  make migrate-up     - Run database migrations up"
	@echo "  make migrate-down   - Run database migrations down"
	@echo "  make kafka-up       - Start Kafka cluster"
	@echo "  make kafka-down     - Stop Kafka cluster"
	@echo "  make kafka-topics   - List Kafka topics"
	@echo "  make kafka-create-topic - Create a new Kafka topic"
	@echo "  make minio-up       - Start MinIO server"
	@echo "  make minio-down     - Stop MinIO server"
	@echo "  make minio-mc       - Start MinIO client shell"
