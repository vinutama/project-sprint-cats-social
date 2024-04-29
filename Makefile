build-dev:
	docker compose build

restart-dev:
	docker restart cats-web

run-dev:
	docker compose up -d

logs-web:
	docker logs -f cats-web

logs-db:
	docker logs -f cats-db

check-db:
	docker exec -it cats-db psql -U cats -d cats-db

clear-db:
	docker rm -f -v cats-db
