package db

import (
	"backend/src/common"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func Migrate(db *gorm.DB) error {
	//db, err := ConnectToPostgres()
	//if err != nil {
	//	log.Fatal(err)
	//}
	err := db.AutoMigrate(&common.User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&common.Keystone{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&common.Reflection{})
	if err != nil {
		return err
	}
	return err
}

func ConnectToPostgres() (db *gorm.DB, err error) {
	host := os.Getenv("POSTGRES_HOST") // The service name of the PostgreSQL container defined in the docker-compose.yml file
	port := 5432                       // Default PostgreSQL port
	user := os.Getenv("POSTGRES_USER")
	dbname := "test_db"
	password := os.Getenv("POSTGRES_PASSWORD")
	sslmode := "disable" // or "require" if SSL is enabled

	// Create connection string
	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, dbname, password, sslmode)

	// Open a connection to the database
	db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return
}
