
up:
	docker-compose up

down:
	docker-compsoe down

test:
	go test ./...

swag:
	swag init -o ./registry/docs -d ./registry -g ./cmd/main.go
	swag init -o ./wallet/docs   -d ./wallet   -g ./cmd/main.go
	swag init -o ./storage/docs  -d ./storage  -g ./cmd/main.go
