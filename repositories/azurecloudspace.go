package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/katasec/ark/resources/v0"
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
	return &AzureCloudSpaceRepository{
		db:        db,
		tableName: GetTableName[AzureCloudSpaceRepository](),
	}
}

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

func (acs *AzureCloudSpaceRepository) CreateCloudSpace(cs resources.AzureCloudspace) (resources.AzureCloudspace, error) {

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
	return cs, nil
}

func (acs *AzureCloudSpaceRepository) UpdateCloudSpace(cs resources.AzureCloudspace) (resources.AzureCloudspace, error) {

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

func (acs *AzureCloudSpaceRepository) DeleteCloudSpace(cs resources.AzureCloudspace) (resources.AzureCloudspace, error) {

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
