package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB is the global database handle. It can be nil if connection fails or is skipped.
var DB *gorm.DB

// InitDB tries to initialize a GORM MySQL connection using env variables.
// If any required variable is missing or connection fails, it logs and leaves DB as nil.
func InitDB() {
	// Require DB credentials to be explicitly provided in environment or .env
	user := GetEnv("DB_USER", "root")
	pass := GetEnv("DB_PASSWORD", "Kiro#2003")
	host := GetEnv("DB_HOST", "127.0.0.1")
	port := GetEnv("DB_PORT", "3306")
	name := GetEnv("DB_NAME", "EventPlanner")

	if user == "" || pass == "" {
		log.Fatalf("database credentials not set. Please create a .env file with DB_USER and DB_PASSWORD and restart the app")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	DB = db
	log.Printf("database connection initialized")
}
