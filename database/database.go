package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// Initialize MySQL driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Conn *gorm.DB

// Init initializes the database connection
func Init(driver string, user string, password string, host string, port string, database string) error {
	var dbURI string
	if driver == "mysql" {
		dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, database)
	} else {
		dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
	}
	db, err := gorm.Open(driver, dbURI)
	if err != nil {
		return err
	}
	Conn = db
	return nil
}
