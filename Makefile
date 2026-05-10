.PHONY: server-up ngrok-up all-start

server-up:
	cd server && env $(grep -v '^#' .env | xargs) go run ./cmd/server

all-start:
	docker compose up

ngrok-up:
	npx ngrok http 8080
