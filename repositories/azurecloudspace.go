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

// func (acs *AzureCloudSpaceRepository) GetCloudSpace(id int) (CloudSpace, error) {
// }

func (acs *AzureCloudSpaceRepository) CreateCloudSpace(cs messages.AzureCloudspace) (messages.AzureCloudspace, error) {

	acsJson, err := json.Marshal(cs)
	if err != nil {
		fmt.Println(err.Error())
	}
	sqlCmd := fmt.Sprintf("insert into cloudspaces(name, cloudspace) values('%s', '%s')", cs.Name, acsJson)
	_, err = acs.db.Exec(sqlCmd)
	if err != nil {
		fmt.Println(err.Error())
	}
	return cs, nil
}
