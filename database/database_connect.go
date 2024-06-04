package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(dsn string, db *gorm.DB) (*gorm.DB, error) {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	fmt.Println("Successfully connected to database!")

	return db, nil
}

// Function to drop and recreate the database
func RecreateDatabase(adminDsn, dbName string, db *gorm.DB) error {
	// Connect to the PostgreSQL server without specifying a database (default is "postgres")
	db, err := ConnectDatabase(adminDsn, db)
	if err != nil {
		return err
	}

	// Terminate active connections to the target database
	if err := terminateConnections(db, dbName); err != nil {
		return fmt.Errorf("failed to terminate active connections: %w", err)
	}

	// Drop the existing database
	if err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)).Error; err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	// Create a new database
	if err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}

// Function to terminate active connections to a database
func terminateConnections(db *gorm.DB, dbName string) error {
	query := fmt.Sprintf(`
        SELECT pg_terminate_backend(pg_stat_activity.pid)
        FROM pg_stat_activity
        WHERE pg_stat_activity.datname = '%s'
        AND pid <> pg_backend_pid();
    `, dbName)
	return db.Exec(query).Error
}
