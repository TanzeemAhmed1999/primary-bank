.PHONY: postgres createdb dropdb migrateup migratedown sqlc

postgres:
	docker run --name primarybank -p 5432:5432 -e POSTGRES_USER=root -e  POSTGRES_PASSWORD=primarybankcode -d postgres:16-alpine

createdb:
	docker exec -it primarybank createdb --username=root --owner=root primarybank

dropdb:
	docker exec -it primarybank dropdb primarybank

migrateup:
	migrate -path db/migrations -database "postgresql://root:primarybankcode@localhost:5432/primarybank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:primarybankcode@localhost:5432/primarybank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

update-repos:
	go mod tidy