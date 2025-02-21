package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/primarybank/api"
	"github.com/primarybank/config"
	db "github.com/primarybank/db/sqlc"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := pgxpool.New(context.Background(), cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	if err = server.Start(cfg.ServerAddr); err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
