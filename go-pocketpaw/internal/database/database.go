package database

import (
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Init initializes the SQLite database connection
func Init(dbPath string) error {
	var initErr error
	once.Do(func() {
		var err error
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			initErr = err
			return
		}
	})
	return initErr
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}

// Close closes the database connection
func Close() error {
	if db == nil {
		return nil
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	
	return sqlDB.Close()
}
