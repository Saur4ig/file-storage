.PHONY: run, build, up, down, re, migrate

build:
	docker build -t storage .

up:
	docker-compose up -d storage

down:
	docker-compose rm -sf storage

run:
	make build
	make up

re:
	make down
	make run

migrate:
	migrate -path internal/database/migrations/ -database "postgresql://admin:adminpass@localhost:5432/filestore?sslmode=disable" -verbose up