package repositories

import (
	"reflect"
	"strings"

	"github.com/katasec/ark/resources"
)

// GetTableName Gets the name of the table from the repository name
func GetTableName[T any]() string {
	var t T

	// Get the name of the repository w/o the package name
	tableName := reflect.TypeOf(t).String()
	tableName = tableName[strings.LastIndex(tableName, ".")+1:]

	// Remove the word "Repository" from the name and make it lowercase
	tableName = strings.Replace(tableName, "Repository", "", -1)
	tableName = strings.ToLower(tableName)

	return tableName
}

type Repositories interface {
	Save(resource resources.Resource) error
	Remove(resource resources.Resource) error
}
