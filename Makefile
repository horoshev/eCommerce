
up:
	docker-compose up -d

down:
	docker-compose stop

test:
	go test ./...

swag:
	swag init -o ./registry/docs -d ./registry -g ./cmd/main.go
