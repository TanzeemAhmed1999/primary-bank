package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/primarybank/api"
	db "github.com/primarybank/db/sqlc"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://root:primarybankcode@localhost:5432/primarybank?sslmode=disable"
	serverAddr = "0.0.0.0:8080"
)

func main() {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.Start(serverAddr); err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
