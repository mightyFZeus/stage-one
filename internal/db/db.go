package db

import (
	"log"
	"time"

	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(connStr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxIdleTime(duration)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)

	if err := migrate(gormDB); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return gormDB, nil
}

func migrate(db *gorm.DB) error {
	statements := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
	}

	// Execute other SQL statements
	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			log.Printf("Error executing statement: %s\n%v\n", stmt, err)
			return err
		}
	}

	return nil
}
