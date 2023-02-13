package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	driver = "postgres"
	url    = "postgres://root:password@localhost/bank-go?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	//connectiong to the databae
	testDB, err = sql.Open(driver, url)

	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
