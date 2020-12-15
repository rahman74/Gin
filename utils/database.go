package utils

import (
	"fmt"
	"log"
	"test/model"

	"github.com/jmoiron/sqlx"

	// Postgres dialect
	_ "github.com/lib/pq"
)

// ConnectDB to get all needed db connections for application
func ConnectDB(config *model.Config) *sqlx.DB {
	return getDBConnection(config)
}

func getDBConnection(config *model.Config) *sqlx.DB {
	var dbConnectionStr string

	dbConnectionStr = fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		config.DatabaseConfig.Host,
		config.DatabaseConfig.Port,
		config.DatabaseConfig.DbName,
		config.DatabaseConfig.Username,
		config.DatabaseConfig.Password,
	)

	db, connError := sqlx.Open("postgres", dbConnectionStr)
	if connError != nil {
		log.Panicln("Error establishing connection to database", connError)
		panic(connError)
	}

	pingError := db.Ping()
	if pingError != nil {
		log.Panicln("Error connecting to database", pingError)
		panic(pingError)
	}

	//TODO: experiment with correct values
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(5)

	return db
}
