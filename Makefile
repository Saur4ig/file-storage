.PHONY: run build init up down re migrate test lint

build:
	docker build -t storage .

up:
	docker-compose up -d storage

down:
	docker-compose down

init: build
	docker-compose up -d db
	# Wait until the database is fully ready
	@echo "Waiting for the database to be ready..."
	@until docker exec $$(docker-compose ps -q db) pg_isready -h db -U admin; do sleep 1; done
	make migrate
	make up

migrate:
	migrate -path internal/database/migrations/ \
		-database "postgresql://admin:adminpass@localhost:5432/filestore?sslmode=disable" -verbose up

run: build up

re: down run

test:
	go test -v ./internal/rest

lint:
	golangci-lint run
