package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectToDb() {

	host := os.Getenv("DBHOST")
	user := os.Getenv("DBUSERNAME")
	password := os.Getenv("DBPASS")
	dbname := os.Getenv("DBNAME")
	port := os.Getenv("DBPORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=EST",
		host, user, password, dbname, port,
	)
	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = Db.AutoMigrate(&models.User{}, &models.Session{}, &models.Service{})
	if err != nil {
		panic("failed to connect database")
	}
}
