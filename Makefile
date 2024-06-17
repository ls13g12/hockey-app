dev: 
	docker compose -f docker-compose-dev.yml up -d
	
mongo:
	docker compose -f docker-compose-dev.yml up -d mongo

stop:
	docker compose -f docker-compose-dev.yml down
