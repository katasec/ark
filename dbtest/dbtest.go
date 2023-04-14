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

	//repo.DropTable(db)
	repo.CreateTable(db)
	acs := genCloudSpace()

	acs.Hub.Name = "test2"

	repo.UpdateCloudSpace(acs)
	repo.DeleteCloudSpace(acs)
}

func OpenDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		fmt.Println(err)
	}

	return db
}
