package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/muling3/bank-go/api"
	db "github.com/muling3/bank-go/db/sqlc"
	"github.com/muling3/bank-go/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load config")
		os.Exit(1)
	}
	
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(*store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}