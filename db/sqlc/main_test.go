package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/primarybank/config"
)

var testStore *Store

func TestMain(m *testing.M) {
	cfg, err := config.Load("../..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	testDB, err := pgxpool.New(context.Background(), cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testStore = NewStore(testDB)
	os.Exit(m.Run())
}
