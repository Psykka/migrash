package main

import (
	_ "database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/spf13/cobra"

	"migrash/cmd"
)

var version = "1.0.0"

func main() {
	fmt.Printf("Migrash v%s\n", version)
	cmd.Execute()
}
