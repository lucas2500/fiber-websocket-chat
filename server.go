package main

import (
	"log"

	"fiber-websocket-chat/routes"
	"fiber-websocket-chat/services"

	"github.com/gofiber/fiber/v2"
)

var ()

func main() {

	app := fiber.New()

	routes.SetupRoutes(app)

	go services.RunChat()

	log.Fatal(app.Listen(":3000"))
}
