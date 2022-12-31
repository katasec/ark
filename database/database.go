package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	HOST     = "localhost"
	DATABASE = "ark"
	USER     = "postgres"
	PASSWORD = "postgres"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func SomeStuff() {
	// Initialize connection string.
	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE)

	// Initialize connection object.
	db, err := sql.Open("postgres", connectionString)
	checkError(err)

	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database")

	// Create table script
	sqlCloudspace := `CREATE TABLE IF NOT EXISTS public."cloudspace"
	(
		id serial,
		cloudspacename VARCHAR(50) NOT NULL,
		tags VARCHAR(50),
		projectname VARCHAR(50) NOT NULL,
		CONSTRAINT "cloudspace_pkey" PRIMARY KEY (cloudspacename)
	)`

	// Create table.
	_, err = db.Exec(sqlCloudspace)
	checkError(err)
	fmt.Println("Finished creating table")
}
