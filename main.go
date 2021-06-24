package main

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	models "ultimate.com/exercise/Models"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB
var err error

func main() {
	godotenv.Load()
	setupDB()
	SetupRouter()
}

func setupDB() {
	dbUserName := os.Getenv("DBUserName")
	dbName := os.Getenv("DBName")
	dbPassword := os.Getenv("DBPassword")
	dbPort := os.Getenv("DBPort")
	hostName := os.Getenv("HostName")
	sslMode := os.Getenv("SSLMode")

	connectionString :=
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", hostName, dbPort, dbUserName, dbPassword, dbName, sslMode)

	db, err = gorm.Open("postgres", connectionString)

	if err != nil {

		println("err : ", err.Error())
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})

	db.AutoMigrate(&models.Task{})

	db.AutoMigrate(&models.TaskAssign{})
}
