package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/katasec/ark/sdk/v0/messages"
)

type Repositories interface {
	// GetCloudSpace(id int) (CloudSpace, error)
	// GetCloudSpaces() ([]CloudSpace, error)
	// CreateCloudSpace(cloudspace CloudSpace) (CloudSpace, error)
	// UpdateCloudSpace(cloudspace CloudSpace) (CloudSpace, error)
	// DeleteCloudSpace(id int) error
}

type AzureCloudSpaceRepository struct {
	db *sql.DB
}

func NewAzureCloudSpaceRepository(db *sql.DB) *AzureCloudSpaceRepository {
	return &AzureCloudSpaceRepository{
		db: db,
	}
}

func (acs *AzureCloudSpaceRepository) CreateTable(db *sql.DB) {

	// Create table
	sql_table := `
	CREATE TABLE IF NOT EXISTS cloudspaces(
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT UNIQUE,
		data text
	);
	`
	_, err := db.Exec(sql_table)
	if err != nil {
		fmt.Println(err)
	}
}

func (acs *AzureCloudSpaceRepository) DropTable(db *sql.DB) {

	// Create table
	sql_table := `
	DROP TABLE IF EXISTS cloudspaces;
	`
	_, err := db.Exec(sql_table)
	if err != nil {
		fmt.Println(err)
	}
}

func (acs *AzureCloudSpaceRepository) CreateCloudSpace(cs messages.AzureCloudspace) (messages.AzureCloudspace, error) {

	jsonAcs, err := json.Marshal(cs)
	if err != nil {
		fmt.Println(err.Error())
	}
	sqlCmd := `
	INSERT INTO cloudspaces(name, data)
	VALUES('%s', '%s');
	`

	sqlCmd = fmt.Sprintf(sqlCmd, cs.Name, jsonAcs)

	//sqlCmd := fmt.Sprintf("insert into cloudspaces(name, data) values('%s', '%s')", cs.Name, acsJson)
	_, err = acs.db.Exec(sqlCmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	return cs, nil
}

func (acs *AzureCloudSpaceRepository) UpdateCloudSpace(cs messages.AzureCloudspace) (messages.AzureCloudspace, error) {

	acsJson, err := json.Marshal(cs)
	if err != nil {
		fmt.Println(err.Error())
	}
	sqlCmd := `
	UPDATE cloudspaces
	SET data = '%s'
	WHERE name = '%s';
	`
	sqlCmd = fmt.Sprintf(sqlCmd, acsJson, cs.Name)

	_, err = acs.db.Exec(sqlCmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	return cs, nil
}

func (acs *AzureCloudSpaceRepository) DeleteCloudSpace(cs messages.AzureCloudspace) (messages.AzureCloudspace, error) {

	sqlCmd := `
	Delete from cloudspaces
	WHERE name = '%s';
	`
	sqlCmd = fmt.Sprintf(sqlCmd, cs.Name)

	_, err := acs.db.Exec(sqlCmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	return cs, nil
}
