package dbtest

import (
	"database/sql"
	"fmt"

	"github.com/katasec/ark/repositories"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	fmt.Println("This is init")
}

func Start() {
	fmt.Println("This is start!")

	DbStuff()
}

func DbStuff() {
	db := OpenDb()
	defer db.Close()

	repo := repositories.NewAzureCloudSpaceRepository(db)
	acs := genCloudSpace()

	repo.CreateCloudSpace(acs)
}
func CreateTable() {

	db := OpenDb()
	defer db.Close()

	// Create table
	sql_table := `
	CREATE TABLE IF NOT EXISTS cloudspaces(
		id integer NOT NULL PRIMARY KEY,
		name text,
		cloudspace text
	);
	`
	_, err := db.Exec(sql_table)
	if err != nil {
		fmt.Println(err)
	}
}

//NewAzureCloudSpaceRepository

func OpenDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		fmt.Println(err)
	}

	return db
}
