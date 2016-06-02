package datastore

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// Initialize MySQL driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/torrent-viewer/backend/herr"
)

// Conn is the database connection used to store data
var Conn *gorm.DB

// Identifiable represent an entity that can be identified by its unique ID
// All of the models used with this datastore must be identifiable.
type Identifiable interface {
	GetID() int
}

// Init initializes the database connection
func Init(driver string, user string, password string, host string, port string, database string) error {
	var dbURI string
	if driver == "mysql" {
		dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, database)
	} else if driver == "sqlite3" {
		dbURI = database
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

// CountEntities count entities from the datastore with the given constraints
func CountEntities(model interface{}, out interface{}, where interface{}, args ...interface{}) *herr.Error {
	conn := Conn.Model(model)
	if where != nil {
		conn = conn.Where(where, args...)
	}
	if err := conn.Count(out).Error; err != nil {
		return &herr.Error{
			ID:     "database-error",
			Status: "500",
			Title:  "Database Error",
			Detail: err.Error(),
		}
	}
	return nil
}

// FetchEntities fetch entities from the datastore with the given constraints
func FetchEntities(out interface{}, where ...interface{}) *herr.Error {
	if err := Conn.Find(out, where...).Error; err != nil {
		return &herr.Error{
			ID:     "database-error",
			Status: "500",
			Title:  "Database Error",
			Detail: err.Error(),
		}
	}
	return nil
}

// FetchEntities fetch entities from the datastore with the given constraints
func FetchPagedEntities(out interface{}, limit int, offset int, where ...interface{}) *herr.Error {
	if err := Conn.Limit(limit).Offset(offset).Find(out, where...).Error; err != nil {
		return &herr.Error{
			ID:     "database-error",
			Status: "500",
			Title:  "Database Error",
			Detail: err.Error(),
		}
	}	
	return nil
}

// FetchEntity fetch an entity based on its ID
func FetchEntity(out interface{}, id int) *herr.Error {
	d := Conn.First(out, id)
	if d.RecordNotFound() != false {
		err := d.Error
		return &herr.Error{
			ID:     "not-found",
			Status: "404",
			Title:  "Not Found",
			Detail: err.Error(),
		}
	} else if err := d.Error; err != nil {
		return &herr.Error{
			ID:     "database-error",
			Status: "500",
			Title:  "Database Error",
			Detail: err.Error(),
		}
	}
	return nil
}

// StoreEntity store a new entity in the datastore.
// The stored entity is not allowed to specify an ID.
func StoreEntity(in interface{}) *herr.Error {
	if Conn.NewRecord(in) != true {
		return &herr.DuplicateEntryError;
	}
	if err := Conn.Create(in).Error; err != nil {
		return &herr.Error{
			ID:     "database-error",
			Status: "500",
			Title:  "Database Error",
			Detail: err.Error(),
		}
	}
	return nil
}

// UpdateEntity update an entity in the datastore
func UpdateEntity(in interface{}) *herr.Error {
	if err := Conn.Model(in).Update(in).Error; err != nil {
		return &herr.Error{
			ID:     "database-error",
			Status: "500",
			Title:  "Database Error",
			Detail: err.Error(),
		}
	}
	return nil
}

// DeleteEntity delete an entity in the datastore,
// using the ID property of the given model.
func DeleteEntity(in Identifiable) *herr.Error {
	var count int
	if err := Conn.Model(in).Where("id = ?", in.GetID()).Count(&count).Error; err != nil {
		return &herr.Error{
			ID:     "database-error",
			Status: "500",
			Title:  "Database Error",
			Detail: err.Error(),
		}
	}
	if count == 0 {
		return &herr.Error{
			ID:     "not-found",
			Status: "404",
			Title:  "Not Found",
			Detail: "The requested resource was not found in the datastore.",
		}
	}
	if err := Conn.Delete(in).Error; err != nil {
		return &herr.Error{
			ID:     "database-error",
			Status: "500",
			Title:  "Database Error",
			Detail: err.Error(),
		}
	}
	return nil
}