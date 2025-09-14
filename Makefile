.PHONY: run migrate-up migrate-down tidy

run:
	go run ./cmd/app serve

migrate-up:
	go run ./cmd/app migrate up

migrate-down:
	go run ./cmd/app migrate down

tidy:
	go mod tidy