# Migrash

Migrash is a simple migration tool for SQL databases. It is written in shell script and is designed to be simple and easy to use in any environment.

## List of supported databases

- [x] PostgreSQL (not tested)
- [x] MySQL/MariaDB (not tested)
- [x] SQLite

## Usage

The following commands are available:

- `migrash create` - Create a new migration file and a directory if it does not exist.
- `migrash up` - Apply all migrations that have not yet been applied.
- `migrash down` - Rollback the last migration.
- `migrash redo` - Rollback the last migration and then apply it again.
- `migrash reset` - Rollback all migrations.
- `migrash status` - Show the status of all migrations.

## Installation

### Unix systems

Install with curl or wget, then run one of the fallowing commands:
```sh
# Install with curl
curl -sSL https://raw.githubusercontent.com/0x111/migrash/master/install.sh | sh

# install with wget
wget -qO- https://raw.githubusercontent.com/0x111/migrash/master/install.sh | sh
```

### Windows

You can install it on Git Bash or another Unix emulator on Windows.

## Configuration

Migrash uses default driver to connect to the database. You can specify the driver in the configuration file. The configuration file is located in the root directory of the project and is called `.migrashrc`. The configuration file is a shell script that is executed before each command. You can use any shell commands in the configuration file.

Example configuration file can be found in the root directory of the project. or [here](./.migrashrc).

## Migration files

Migration files are located in the `migrations` directory. The name of the migration dir will be in the format `YYYYMMDDHHMMSS-name`, inside the directory there will be two files `up.sql` and `down.sql`. The `up.sql` file contains the SQL query to apply the migration, and the `down.sql` file contains the SQL query to rollback the migration.

## License

This project is licensed under the terms of the [WTFPL](./LICENSE) license.

