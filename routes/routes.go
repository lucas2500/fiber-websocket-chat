package routes

import (
	"fiber-websocket-chat/middleware"
	"fiber-websocket-chat/services"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/ws/NewChat", middleware.IsWebSocketUpgrade, websocket.New(services.NewChat))
}
