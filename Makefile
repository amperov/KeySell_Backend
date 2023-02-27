SHELL := /bin/bash

DB_VERSION := 15
DB_USER := unflat
DB_PASSWORD := 123
DB_NAME := sell
DB_PORT := 5432
DB_HOST := localhost


all_db: db-download db-run migrate-up
db-download:
	echo "Pulling Container"
	docker pull postgres:$(DB_VERSION)

db-run:
	echo "Running docker container"
	docker run --name=$(DB_NAME) \
	-e POSTGRES_USER=$(DB_USER) \
	-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
	-e POSTGRES_DB=$(DB_NAME) \
 	-p $(DB_PORT):$(DB_PORT) -d --rm postgres:$(DB_VERSION)

migrate-up:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down
