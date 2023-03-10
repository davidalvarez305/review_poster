package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DBInstance = *gorm.DB

var DB DBInstance

type connection struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

func Connect() {
	conn := connection{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DB:       os.Getenv("POSTGRES_DB"),
	}

	db, err := gorm.Open(postgres.Open(connToString(conn)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		FullSaveAssociations: true,
	})

	if err != nil {
		fmt.Printf("Error connecting to the DB: %s\n", err.Error())
		return
	} else {
		fmt.Printf("Connected to Database.")
	}

	DB = db
}

func connToString(info connection) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		info.Host, info.Port, info.User, info.Password, info.DB)
}
