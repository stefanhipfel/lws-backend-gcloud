package common

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var session *sql.DB

var createTableStatements = []string{
	`CREATE DATABASE IF NOT EXISTS website DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`,
	`USE website;`,
	`CREATE TABLE IF NOT EXISTS blogs (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		title VARCHAR(255) NULL,
		category VARCHAR(255) NULL,
		createdOn DATE NULL,
		imageUrl VARCHAR(255) NULL,
		text TEXT NULL,
		createdBy VARCHAR(255) NULL,
		PRIMARY KEY (id)
	)`,
}

func GetMySQLSession() *sql.DB {
	if session != nil {
		return session
	}
	createDbSession()
	return session
}

func createDbSession() {
	var err error
	connectionName := mustGetenv("CLOUDSQL_CONNECTION_NAME")
	user := mustGetenv("CLOUDSQL_USER")
	password := os.Getenv("CLOUDSQL_PASSWORD") // NOTE: password may be empty

	session, err = sql.Open("mysql", fmt.Sprintf("%s:%s@cloudsql(%s)/website", user, password, connectionName))
	if err != nil {
		fmt.Println(err.Error())
		log.Fatalf("[createDbSession]: %s\n", err)
	}

	err = createTable(session)
	if err != nil {
		log.Fatalf("Could not query db: %v", err)
		fmt.Println(err.Error())
		return
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}

// createTable creates the table, and if necessary, the database.
func createTable(conn *sql.DB) error {
	for _, stmt := range createTableStatements {
		_, err := conn.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}
