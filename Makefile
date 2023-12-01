server:
	go run cmd/bot/main.go

docker-up:
	docker compose -f docker-compose.yaml up --build

docker-down:
	docker compose down