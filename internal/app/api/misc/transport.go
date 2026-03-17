package misc

import "github.com/gofiber/fiber/v3"

func VersionHandler(service Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		resp, err := service.Version(c.Context())
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return err
		}

		return c.JSON(resp)
	}
}

func HealthHandler(service Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		health, err := service.Health(c.Context())
		if err != nil {
			return err
		}

		return c.JSON(health)
	}
}
