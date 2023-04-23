package repositories

import (
	"reflect"
	"strings"
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
