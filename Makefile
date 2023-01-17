
server:
	go run cmd/main.go

build:
	go build -o bin/server cmd/main.go

d.up:
	sudo docker-compose --env-file ./cartique/.env up

d.down:
	sudo docker-compose down

d.up.build:
	docker-compose --build up

d.upd:
	sudo docker-compose --env-file ./cartique/.env up -d
	