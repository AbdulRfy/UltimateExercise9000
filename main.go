package main

import (
	"ultimate.com/exercise/controller"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	controller.SetupDB()
	controller.SetupRouter()
}
