.PHONY: postgres createdb dropdb migrateup migratedown sqlc mockgen update-repos test server

postgres:
	docker run --name primarybank --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e  POSTGRES_PASSWORD=primarybankcode -d postgres:16-alpine

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

mockgen:
	mockgen -destination=db/mocks/store.go -package=mocks github.com/primarybank/db/sqlc Store

update-repos:
	go mod tidy

test:
	go test -v -cover ./...

server:
	go run main.go