.PHONY: launch_services stop_services build_services

launch_services: build_services migrate
	docker compose up

stop_services:
	docker compose --profile test down

build_services:
	docker compose build --no-cache

migrate:
	docker compose run --rm migrate