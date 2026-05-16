package auth

import "github.com/gofiber/fiber/v3"

func Route(app fiber.Router, svc Service) {
	app.Post("/", CreateTokenHandler(svc))
	app.Delete("/", Middleware(svc), RevokeTokenHandler(svc))
}
