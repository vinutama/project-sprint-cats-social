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

migrate-db:
	migrate -database "postgres://cats:password@localhost:5432/cats-db?sslmode=disable" -path db/migrations up

migrate-db-down:
	migrate -database "postgres://cats:password@localhost:5432/cats-db?sslmode=disable" -path db/migrations down -all

build-prod-linux:
	GOOS=linux GOARCH=amd64 go build -o build/main_kangman53

build-prod-win:
	GOOS=windows GOARCH=amd64 go build -o build/cat-social-win.exe

build-prod-mac:
	GOOS=darwin GOARCH=amd64 go build -o build/main_kangman53
