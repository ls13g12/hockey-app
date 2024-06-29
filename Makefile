include .env

dev: 
	docker compose -f docker-compose-dev.yml up -d
	
mongo:
	docker compose -f docker-compose-dev.yml up -d mongo

tui-dev: mongo
	cd ./src/cmd/tui && go run main.go -mode=dev

tui-prod: 
	cd ./src/cmd/tui && go run main.go -mode=prod -dsn=${MONGODB_URI}

stop:
	docker compose -f docker-compose-dev.yml down

remove-all:
	docker container rm $$(docker ps -aq) -f
