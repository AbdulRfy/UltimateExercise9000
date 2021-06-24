package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB
var err error

func main() {
	godotenv.Load()
	setupDB()
	setupRouter()
}

func setupDB() {
	dbUserName := os.Getenv("DBUserName")
	dbName := os.Getenv("DBName")
	dbPassword := os.Getenv("DBPassword")

	connectionString :=
		fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable", dbUserName, dbPassword, dbName)

	db, err = gorm.Open("postgres", connectionString)

	if err != nil {

		println("err : ", err.Error())
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	db.AutoMigrate(&Task{})

	db.AutoMigrate(&TaskAssign{})
}


