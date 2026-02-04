.PHONY: compose-up up down logs swag

compose-up up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f

swag:
	swag init -g cmd/api/main.go