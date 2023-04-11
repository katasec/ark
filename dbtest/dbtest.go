package dbtest

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	fmt.Println("This is init")
}

func Start() {
	fmt.Println("This is start!")

	acs := genCloudSpace()
	fmt.Println(acs)
	DbStuff()
}

func DbStuff() {

	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Create table
	sql_table := `
	CREATE TABLE IF NOT EXISTS cloudspaces(
		id integer NOT NULL PRIMARY KEY,
		name text,
		cloudspace text
	);
	`
	_, err = db.Exec(sql_table)
}
