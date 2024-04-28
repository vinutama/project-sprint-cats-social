build-dev:
	docker compose build

run-dev:
	docker compose up -d

logs-web:
	docker logs -f cats-web

logs-db:
	docker logs -f cats-db

check-db:
	docker exec -it cats-db psql -U cats -d cats-db