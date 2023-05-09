package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/katasec/ark/resources/v0/azure/cloudspaces"
)

type Repositories interface {
	// GetCloudSpace(id int) (CloudSpace, error)
	// GetCloudSpaces() ([]CloudSpace, error)
	// CreateCloudSpace(cloudspace CloudSpace) (CloudSpace, error)
	// UpdateCloudSpace(cloudspace CloudSpace) (CloudSpace, error)
	// DeleteCloudSpace(id int) error
}

type AzureCloudSpaceRepository struct {
	db        *sql.DB
	tableName string
}

func NewAzureCloudSpaceRepository(db *sql.DB) *AzureCloudSpaceRepository {

	// Create New Cloud Space Repository
	acsrepo := &AzureCloudSpaceRepository{
		db:        db,
		tableName: GetTableName[AzureCloudSpaceRepository](),
	}

	// Create the table if it doesn't exist
	acsrepo.CreateTable(db)

	return acsrepo
}

// CreateTable creates the table for the repository in the db that's passed in
func (acs *AzureCloudSpaceRepository) CreateTable(db *sql.DB) {

	// Create table
	sql_table := `
	CREATE TABLE IF NOT EXISTS %s(
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT UNIQUE,
		data text
	);
	`
	sql_table = fmt.Sprintf(sql_table, acs.tableName)

	_, err := db.Exec(sql_table)
	if err != nil {
		fmt.Println(err)
	}
}

func (acs *AzureCloudSpaceRepository) DropTable(db *sql.DB) {

	// Create table
	sql_table := `
	DROP TABLE IF EXISTS %s;
	`

	sql_table = fmt.Sprintf(sql_table, acs.tableName)

	_, err := db.Exec(sql_table)
	if err != nil {
		fmt.Println(err)
	}
}

func (acs *AzureCloudSpaceRepository) CreateCloudSpace(cs *cloudspaces.AzureCloudspace) (cloudspaces.AzureCloudspace, error) {

	jsonAcs, err := json.Marshal(cs)
	if err != nil {
		fmt.Println(err.Error())
	}
	sqlCmd := `
	INSERT INTO %s(name, data)
	VALUES('%s', '%s');
	`
	sqlCmd = fmt.Sprintf(sqlCmd, acs.tableName, cs.Name, jsonAcs)
	_, err = acs.db.Exec(sqlCmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	return *cs, nil
}

func (acs *AzureCloudSpaceRepository) UpdateCloudSpace(cs cloudspaces.AzureCloudspace) (cloudspaces.AzureCloudspace, error) {

	jsonData, err := json.Marshal(cs)
	if err != nil {
		fmt.Println(err.Error())
	}
	sqlCmd := `
	UPDATE %s
	SET data = '%s'
	WHERE name = '%s';
	`
	sqlCmd = fmt.Sprintf(sqlCmd, acs.tableName, string(jsonData), cs.Name)

	_, err = acs.db.Exec(sqlCmd)

	if err != nil {
		fmt.Println(err.Error())
	}
	return cs, nil
}

func (acs *AzureCloudSpaceRepository) DeleteCloudSpace(cs cloudspaces.AzureCloudspace) (cloudspaces.AzureCloudspace, error) {

	sqlCmd := `
	Delete from %s
	WHERE name = '%s';
	`
	sqlCmd = fmt.Sprintf(sqlCmd, acs.tableName, cs.Name)

	_, err := acs.db.Exec(sqlCmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	return cs, nil
}

func (acs *AzureCloudSpaceRepository) GetCloudSpaces() (acss []cloudspaces.AzureCloudspace, err error) {

	sqlCmd := `select * from %s`
	sqlCmd = fmt.Sprintf(sqlCmd, acs.tableName)

	fmt.Println(sqlCmd)
	rows, err := acs.db.Query(sqlCmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var data string
		var cs cloudspaces.AzureCloudspace

		err := rows.Scan(&id, &name, &data)
		if err != nil {
			fmt.Println(err)
		}

		err = json.Unmarshal([]byte(data), &cs)
		if err != nil {
			fmt.Println(err)
		}

		acss = append(acss, cs)
	}

	return acss, nil
}

func (acs *AzureCloudSpaceRepository) GetCloudSpace(name string) (cloudspaces.AzureCloudspace, error) {
	var cs cloudspaces.AzureCloudspace

	sqlCmd := `select * from %s where name='%s'`
	sqlCmd = fmt.Sprintf(sqlCmd, acs.tableName, name)

	fmt.Println(sqlCmd)
	rows, err := acs.db.Query(sqlCmd)
	if err != nil {
		log.Println("GetCloudSpace:", err.Error())
		return cloudspaces.AzureCloudspace{}, err
	}
	defer rows.Close()

	// Create empty cloudspace if no rows are returned
	if !rows.Next() {
		fmt.Println("No rows returned")
		return *cloudspaces.NewAzureCloudSpace(), nil
	}

	fmt.Println("Found something")
	// Else return the cloudspace
	var data string
	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name, &data)
		if err != nil {
			fmt.Println(err)
		}

		err = json.Unmarshal([]byte(data), &cs)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("The data is: ", data)
	}
	fmt.Println("The data is: ", data)

	return cs, nil
}
