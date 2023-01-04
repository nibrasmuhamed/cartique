package main

import (
	"log"

	"github.com/nibrasmuhamed/cartique/database"
	"github.com/nibrasmuhamed/cartique/routes"
)

func main() {
	database.InitDB()
	app := routes.Routes()
	log.Fatal(app.Listen(":8000"))
}
