package database

import (
	"fmt"

	"migrash/internal/config"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx"
	_ "github.com/mattn/go-sqlite3"
)

func getConnection(dbms string, url string) *sqlx.DB {
	var db *sqlx.DB
	var err error

	switch dbms {
	case "mysql", "mariadb":
		db, err = sqlx.Open("mysql", url)
	case "postgres":
		db, err = sqlx.Open("pgx", url)
	case "sqlite3":
		db, err = sqlx.Open("sqlite3", url)
	}

	if err != nil {
		panic(err)
	}

	return db
}

func Connect(config *config.Config) *sqlx.DB {
	databaseVersion := fmt.Sprintf("Starting a connection with %s on server %s@%s:%s", config.Database.DBMS, config.Database.User, config.Database.Host, config.Database.Port)
	fmt.Println(databaseVersion)

	db := getConnection(config.Database.DBMS, config.Database.Url)

	// TODO: pass the database name sanitized
	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + config.Database.Name); err != nil {
		panic(err)
	}

	db = getConnection(config.Database.DBMS, config.Database.Url+config.Database.Name)
	fmt.Println("Connected to", config.Database.DBMS)

	return db
}
