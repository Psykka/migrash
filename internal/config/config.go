package config

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Database struct {
	DBMS       string
	Url        string
	Host       string
	Port       string
	User       string
	Password   string
	Name       string
	QueryParam map[string]string
}

type Config struct {
	MigrationDir   string
	MigrationTable string
	Database       Database
}

var databases = []string{
	"sqlite3",
	"mysql",
	"postgres",
}

var configFileName = ".migrashrc"

func parseConnectionString(connString string, database *Database) error {
	isSqlite := isSQLite(connString)
	if isSqlite {
		database.DBMS = "sqlite3"
		return nil
	}

	parsedURL, err := url.Parse(connString)
	if err != nil {
		return err
	}

	if database.Name == "" {
		return fmt.Errorf("unsupported connection string: database name is not set")
	}

	database.DBMS = getDBMS(parsedURL.Scheme)

	if database.DBMS == "" {
		return fmt.Errorf("unsupported connection string: scheme must be one of %v", strings.Join(databases, ", "))
	}

	user := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()
	host, port := parseHostAndPort(parsedURL.Host)
	name := strings.TrimLeft(parsedURL.Path, "/")
	queryParams := parseQueryParams(parsedURL.RawQuery)

	database.setURL(user, password, host, port)
	database.Host = host
	database.Port = port
	database.User = user
	database.Password = password
	database.Name = name
	database.QueryParam = queryParams

	return nil
}

func isSQLite(connString string) bool {
	pattern := `\.(sqlite|sqlite3|db|db3)\b`
	re := regexp.MustCompile(pattern)
	return re.MatchString(connString)
}

func getDBMS(scheme string) string {
	for _, prefix := range databases {
		if strings.HasPrefix(scheme, prefix) {
			return prefix
		}
	}
	return ""
}

func parseHostAndPort(hostPort string) (string, string) {
	parts := strings.Split(hostPort, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

func parseQueryParams(rawQuery string) map[string]string {
	queryParams := make(map[string]string)
	if len(rawQuery) > 0 {
		queryString := strings.Split(rawQuery, "&")
		for _, param := range queryString {
			pair := strings.Split(param, "=")
			if len(pair) == 2 {
				queryParams[pair[0]] = pair[1]
			}
		}
	}
	return queryParams
}

func (database *Database) setURL(user, password, host, port string) {
	if database.DBMS == "postgres" {
		database.Url = fmt.Sprintf("postgres://%s:%s@%s:%s/", user, password, host, port)
	} else {
		database.Url = fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, password, host, port)
	}
}

func parseConfigFile(filename string) (*Config, error) {
	configFile, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer configFile.Close()

	config := &Config{}

	scanner := bufio.NewScanner(configFile)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "#") && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.Trim(strings.TrimSpace(parts[1]), "\"")

			if strings.Contains(value, "$") {
				value, _ = os.LookupEnv(strings.Trim(value, "$"))
			}

			switch key {
			case "MIGRATION_DIR":
				{
					config.MigrationDir = value
				}
			case "DATABASE_URL":
				{
					config.Database.Url = value
					err := parseConnectionString(value, &config.Database)

					if err != nil {
						return nil, err
					}
				}
			case "MIGRATION_TABLE":
				{
					config.MigrationTable = value
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if config.Database.Url == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	if config.MigrationDir == "" {
		config.MigrationDir = "migrations"
	}

	if config.MigrationTable == "" {
		config.MigrationTable = "_migrash_migrations"
	}

	return config, nil
}

func parseConfigEnv() (*Config, error) {
	config := &Config{}

	migrationDir, ok := os.LookupEnv("MIGRATION_DIR")
	if !ok {
		migrationDir = "migrations"
	}

	config.MigrationDir = migrationDir

	migrationTable, ok := os.LookupEnv("MIGRATION_TABLE")
	if !ok {
		migrationTable = "_migrash_migrations"
	}

	config.MigrationTable = migrationTable

	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	config.Database.Url = databaseUrl
	err := parseConnectionString(databaseUrl, &config.Database)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func ParseConfig() (*Config, error) {
	_, err := os.Stat(configFileName)

	if err != nil {
		return parseConfigEnv()
	}

	return parseConfigFile(configFileName)
}
