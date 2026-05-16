package auth

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/lvjp/wtf-go/pkg/api"
)

func CreateTokenHandler(svc Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		var req api.AuthTokenCreateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return err //nolint:wrapcheck // Let Fiber handle the error response.
		}

		tok, err := svc.CreateToken(c.Context(), req.Subject)
		if err != nil {
			if errors.Is(err, errAuthenticationFailed) {
				return errors.Join(fiber.ErrForbidden, err)
			}
			return err
		}

		c.Status(fiber.StatusCreated)
		return c.JSON(api.AuthTokenResponse{
			ID:       tok.ID,
			NotAfter: tok.NotAfter,
		})
	}
}

func RevokeTokenHandler(svc Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := TokenFromContext(c)

		if err := svc.RevokeToken(c.Context(), token.ID); err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
