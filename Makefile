.PHONY: run, build, up, down, re

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