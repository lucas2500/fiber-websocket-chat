package main

import (
	"log"
	"sync"

	"fiber-websocket-chat/routes"
	"fiber-websocket-chat/services"

	"github.com/gofiber/fiber/v2"
)

func main() {

	var wg sync.WaitGroup

	app := fiber.New()
	routes.SetupRoutes(app)

	go services.RunChat()

	log.Fatal(app.Listen(":3000"))

	wg.Wait()
}
