package database

import (
	"github.com/jinzhu/gorm"
	// Initialize MySQL driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Conn *gorm.DB

// Init initializes the database connection
func Init(driver string, database string) error {
	db, err := gorm.Open(driver, database)
	if err != nil {
		return err
	}
	Conn = db
	return nil
}
