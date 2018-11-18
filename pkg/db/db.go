package db

import "time"
import "github.com/boltdb/bolt"

// CreateTableIfNotExists creates specified table if it doesn't exist
func CreateTableIfNotExists(tableName string) error {
	return nil
}

// DB is the singleton of the database
// it should be InitDB before used
var DB *bolt.DB

// InitDB creates the singleton of database
func InitDB(dbPath string) error {
	var err error
	DB, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	return err
}
