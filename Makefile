postgres:
	docker run --name crawler_db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=pwd123 -d postgres:12-alpine

createdb:
	winpty docker exec -it crawler_db createdb --username=root --owner=root crawler_db

dropdb:
	winpty docker exec -it crawler_db dropdb crawler_db

forcedbversion:
	migrate -path db/migration -database "postgresql://root:pwd123@localhost:5432/crawler_db?sslmode=disable" force 1

migrateup:
	migrate -path db/migration -database "postgresql://root:pwd123@localhost:5432/crawler_db?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:pwd123@localhost:5432/crawler_db?sslmode=disable" -verbose down 1

sqlc:
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb migrateup sqlc