package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupGormClient() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=mini_wallet user=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}
	log.Println("gorm client setup finished")

	return db
}
