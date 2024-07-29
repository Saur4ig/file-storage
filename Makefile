.PHONY: run, build, up, down

build:
	docker build -t storage .

up:
	docker-compose up -d storage

down:
	docker-compose rm -sf storage

run:
	make build
	make up