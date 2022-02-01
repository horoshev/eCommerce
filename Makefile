run:
	docker-compose up -d --build registry storage wallet kafka zookeeper mongodb-primary mongodb-secondary mongodb-arbiter

up:
	docker-compose up -d

down:
	docker-compose stop

test:
	go test ./...

swag:
	swag init -o ./registry/docs -d ./registry -g ./cmd/main.go
