package server

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq" // Import the pq driver
)

func TestServerDbConfig(t *testing.T) {

	server := NewServer()
	db, err := server.GetDbConnection()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

}
