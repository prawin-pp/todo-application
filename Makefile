#!/usr/bin/make

include .env

SHELL = /bin/bash

CURRENT_UID := $$(id -u)
CURRENT_GID := $$(id -g)
WORKDIRS := $$(pwd)

DB_CONN ?= postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable
MIGRATE := docker run --rm --user $(CURRENT_UID):$(CURRENT_GID) -v $(WORKDIRS)/api/migrations:/migrations --network development migrate/migrate -path=/migrations/ -database "$(DB_CONN)"

default: dev ## start dev containers

.PHONY: dev
dev: ## start dev containers
	@echo "Starting dev containers..."
	@docker compose down
	@docker compose -f docker-compose.yml up -d --build

.PHONY: log
log: ## show logs
	@echo "Showing logs..."
	@docker compose -f docker-compose.yml logs -f

.PHONY: migrate-create
migrate-create: ## create a new database migration
	@read -p "Enter the name of the new migration: " name; \
	$(MIGRATE) create -ext sql -seq -dir /migrations/ $${name// /_}

.PHONY: migrate
migrate: ## run all new database migrations
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## revert database to the last migration step
	@echo "Reverting database to the last migration step..."
	@$(MIGRATE) down 1

.PHONY: migrate-drop
migrate-drop: ## drop tables
	@echo "Droping tables..."
	@$(MIGRATE) drop -f

.PHONY: migrate-reset
migrate-reset: ## reset database and re-run all migrations
	@echo "Resetting database..."
	@$(MIGRATE) drop -f
	@echo "Running all database migrations..."
	@$(MIGRATE) up
