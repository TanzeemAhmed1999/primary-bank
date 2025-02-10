package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:primarybankcode@localhost:5432/primarybank?sslmode=disable"
)

var testStore *Store

func TestMain(m *testing.M) {
	var err error
	testDB, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testStore = NewStore(testDB)
	os.Exit(m.Run())
}
