package main

import (
	"golang-todos/database"
	"golang-todos/routes"
)

func main() {
	database.Connect("root:admin@123@tcp(localhost:3306)/golang?parseTime=true")
	database.Migrate()

	router := routes.InitRouter()
	router.Run(":3001") // listen and serve on 0.0.0.0:3001
}
