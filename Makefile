.PHONY: start build/prod start/prod push/prod

start:
	go run ./cmd/server

build/prod:
	docker compose -f docker-compose.prod.yaml build

start/prod:
	docker compose -f docker-compose.prod.yaml up --build

push/prod:
	docker compose -f docker-compose.prod.yaml push

