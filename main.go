package main

import (
	_ "database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/spf13/cobra"

	"migrash/cmd"
)

func main() {
	cmd.Execute()
}
