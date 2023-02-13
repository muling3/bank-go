package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/muling3/bank-go/util"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, errr := util.LoadConfig("../..")
	if errr != nil {
		fmt.Println("cannot load config")
		os.Exit(1)
	}
	
	var err error

	//connectiong to the databae
	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
