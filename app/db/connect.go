package db

import (
	"database/sql"
	"go-gacha/app/config"
	"log"
)

func Connect() *sql.DB {
	log.Printf("Server listening on http://localhost:%s", config.Config.Port)
	db, err := sql.Open(config.Config.DriverName, config.Config.DataSourceName)
	if err != nil {
		panic(err.Error())
	}
	return db
}