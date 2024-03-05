# Migrash

Migrash a language-agnostic simple SQL migration tool. It is designed to be simple and easy to use. It is written in Go and is designed to be used in a CI/CD pipelines.

## List of supported databases

- [ ] PostgreSQL
- [ ] MySQL/MariaDB
- [ ] SQLite

## Usage

The following commands are available:

- `migrash init` - Create config file for migrash
- `migrash create` - Create a new migration file and a directory if it does not exist.
- `migrash up` - Apply all migrations that have not yet been applied.
- `migrash down` - Rollback the last migration.
- `migrash redo` - Rollback the last migration and then apply it again.
- `migrash reset` - Rollback all migrations.
- `migrash status` - Show status check if there are pending or executed migrations.

## Installation

Soon...

## Configuration

Migrash uses DATABASE_URL to connect to the database. The DATABASE_URL environment variable is a URL that specifies the database connection. The URL is in the format `dialect://user:password@host:port/database?param1=value1&param2=value2`.

The configuration file is located in the root directory of the project and is called `.migrashrc`. The configuration file is a text file that used before each command.

Example configuration file can be found in the root directory of the project. or [here](./.migrashrc).

## Migration files

Migration files are located in the `migrations` directory by default. The name of the migration dir will be in the format `YYYYMMDDHHMMSS-name`, inside the directory there will be two files `up.sql` and `down.sql`. The `up.sql` file contains the SQL query to apply the migration, and the `down.sql` file contains the SQL query to rollback the migration.

## License

This project is licensed under the terms of the [WTFPL](./LICENSE) license.
