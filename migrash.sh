#!/bin/sh

# Migrash - A simple sql migration tool in shell script
# Author's Name: Psykka

# This shell script is distributed under the terms of the
# WTFPL (Do What The F*ck You Want To Public License).
#
# This script is provided without any warranty or guarantee of any kind.
# You should have received a copy of the WTFPL along with this script.
# If you didn't, you can obtain a copy from the WTFPL website:
# http://www.wtfpl.net

echo "Migrash v1.0.0"

CURRENT_TIMESTAMP=$(date +'%Y%m%d%H%M%S')

throw_error() {
    echo $1
    exit 1
}

check_error() {
    if [ $? -ne 0 ]; then
        throw_error $1
    fi
}

# Create migrations directory if not exists
create_migration_table() {
    case "$SGDB" in
        sqlite3)
            $SGDB $DATABASE_URL "CREATE TABLE IF NOT EXISTS $MIGRATION_TABLE (id INTEGER PRIMARY KEY, name TEXT NOT NULL, created_at TEXT NOT NULL);"
            ;;
        mysql)
            $SGDB $DATABASE_URL "CREATE TABLE IF NOT EXISTS $MIGRATION_TABLE (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, created_at DATETIME NOT NULL);"
            ;;
        psql)
            $SGDB $DATABASE_URL "CREATE TABLE IF NOT EXISTS $MIGRATION_TABLE (id SERIAL PRIMARY KEY, name VARCHAR(255) NOT NULL, created_at TIMESTAMP NOT NULL);"
            ;;
        *)
            throw_error "SGDB not supported"
            ;;
    esac

    check_error "Error creating migration table"
}

# Insert migration into migrations table
insert_migration() {
    $SGDB $DATABASE_URL "INSERT INTO $MIGRATION_TABLE (name, created_at) VALUES ('$1', '$CURRENT_TIMESTAMP');"
    check_error "Error inserting migration"
}

# Remove migration from migrations table
remove_migration() {
    $SGDB $DATABASE_URL "DELETE FROM $MIGRATION_TABLE WHERE name = '$1';"
    check_error "Error deleting migration"
}

# Get last migration from migrations table
get_last_migration() {
    $SGDB $DATABASE_URL "SELECT * FROM $MIGRATION_TABLE ORDER BY id DESC LIMIT 1;"
    check_error "Error getting last migration"
}

# Check if migration exists in migrations table
migration_exists() {
    $SGDB $DATABASE_URL "SELECT COUNT(*) FROM $MIGRATION_TABLE WHERE name = '$1';"
    check_error "Error checking if migration exists"
}

# Create migration files
create() {
    if [ -z "$1" ]; then
        throw_error "Migration name is required"
    fi

    create_migration_table

    MIGRATION_NAME=$1
    MIGRATION_DIR="$MIGRATIONS_DIR/$CURRENT_TIMESTAMP-$MIGRATION_NAME"
    MIGRATION_UP_FILE="$MIGRATION_DIR/up.sql"
    MIGRATION_DOWN_FILE="$MIGRATION_DIR/down.sql"

    mkdir -p $MIGRATION_DIR
    check_error "Error creating migrations directory"

    touch $MIGRATION_UP_FILE
    check_error "Error creating migration up file"

    touch $MIGRATION_DOWN_FILE
    check_error "Error creating migration down file"

    echo "Migration files created successfully"
}

# Run all pending migrations
up() {
    create_migration_table

    MIGRATIONS=$(ls -1 $MIGRATIONS_DIR | sort -n)
    check_error "Error listing migrations"

    for MIGRATION in $MIGRATIONS; do
        MIGRATION_UP_FILE="$MIGRATIONS_DIR/$MIGRATION/up.sql"

        if [ ! -f "$MIGRATION_UP_FILE" ]; then
            throw_error "Migration up file not found"
        fi

        MIGRATION_EXISTS=$(migration_exists $MIGRATION)

        if [ $MIGRATION_EXISTS -eq 0 ]; then
            echo "Running migration $MIGRATION"

            $SGDB $DATABASE_URL < $MIGRATION_UP_FILE
            check_error "Error running migration"

            insert_migration $MIGRATION
        fi
    done

    echo "Migrations executed successfully"
}

# Rollback last migration
down() {
    create_migration_table

    MIGRATION=$(get_last_migration)

    if [ -z "$MIGRATION" ]; then
        throw_error "No migrations to rollback"
    fi

    MIGRATION_ID=$(echo $MIGRATION | cut -d'|' -f1)
    MIGRATION_NAME=$(echo $MIGRATION | cut -d'|' -f2)
    MIGRATION_DOWN_FILE="$MIGRATIONS_DIR/$MIGRATION_NAME/down.sql"

    if [ ! -f "$MIGRATION_DOWN_FILE" ]; then
        throw_error "Migration down file not found"
    fi

    echo "Rolling back migration $MIGRATION_NAME"

    $SGDB $DATABASE_URL < $MIGRATION_DOWN_FILE
    check_error "Error running migration"

    remove_migration $MIGRATION_NAME

    echo "Migration rolled back successfully"
}

# Rollback last migration and run it again
redo() {
    down
    up
}

# Rollback all migrations
reset() {
    create_migration_table

    MIGRATIONS=$(ls -1 $MIGRATIONS_DIR | sort -nr)
    check_error "Error listing migrations"

    for MIGRATION in $MIGRATIONS; do
        MIGRATION_DOWN_FILE="$MIGRATIONS_DIR/$MIGRATION/down.sql"

        if [ ! -f "$MIGRATION_DOWN_FILE" ]; then
            throw_error "Migration down file not found"
        fi

        echo "Rolling back migration $MIGRATION"

        $SGDB $DATABASE_URL < $MIGRATION_DOWN_FILE
        check_error "Error running migration"

        remove_migration $MIGRATION
    done

    echo "Migrations rolled back successfully"
}

# Show migrations status
status() {
    create_migration_table

    MIGRATIONS=$(ls -1 $MIGRATIONS_DIR | sort -n)
    check_error "Error listing migrations"

    for MIGRATION in $MIGRATIONS; do
        MIGRATION_EXISTS=$(migration_exists $MIGRATION)

        if [ $MIGRATION_EXISTS -eq 0 ]; then
            echo "Pending: $MIGRATION"
        else
            echo "Executed: $MIGRATION"
        fi
    done
}

show_help() {
    echo "Usage: migrash.sh [options]"
    echo "Options:"
    echo "  create        Create a new migration"
    echo "  up            Run all pending migrations"
    echo "  down          Rollback last migration"
    echo "  redo          Rollback last migration and run it again"
    echo "  reset         Rollback all migrations"
    echo "  status        Show migrations status"
    echo "  help          Show this help"
}

# call .rc file
if [ -f "./.migrashrc" ]; then
    . ./.migrashrc
fi

case "$1" in
    create)
        create $2
        ;;
    up)
        up
        ;;
    down)
        down
        ;;
    redo)
        redo
        ;;
    reset)
        reset
        ;;
    status)
        status
        ;;
    help)
        show_help
        ;;
    *)
        show_help
        ;;
esac

exit 0