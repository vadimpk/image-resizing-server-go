.PHONY:
.SILENT:

build-api:
	go build -o ./.bin/api cmd/api/main.go

run-api: build-api
	./.bin/api

build-opt:
	go build -o ./.bin/optimizer cmd/optimizer/main.go

run-opt: build-opt
	./.bin/optimizer

run-rabbit:
	docker run -d --hostname my-rabbit --name rabbit-images -p 15672:15672 -p 5672:5672 rabbitmq:3-management