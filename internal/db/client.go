package db

import (
	"fmt"
	"time"

	"xarr-proxy/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func Init(cfg *config.Config) {
	var err error
	db, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s/%s", cfg.ConfDir, "xarrproxy.db")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	if sqlDB, err := db.DB(); err != nil {
		panic(err)
	} else {
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(10)

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(100)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
}

func Get() *gorm.DB {
	return db
}
