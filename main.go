package main

import (
	"database/sql"
	"github.com/stevenysy/simplebank/util"
	"log"

	_ "github.com/lib/pq"
	"github.com/stevenysy/simplebank/api"
	db "github.com/stevenysy/simplebank/db/sqlc"
)

func main() {
	// Load the configuration
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
}
