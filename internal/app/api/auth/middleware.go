package auth

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"

	"github.com/lvjp/wtf-go/pkg/auth/token"
)

type contextKey struct{}

var localKeyToken = contextKey{}

// Middleware validates the Bearer token from the Authorization header.
// On success, the token is available via TokenFromContext.
func Middleware(svc Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get(fiber.HeaderAuthorization)
		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		id, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found || id == "" {
			return fiber.ErrUnauthorized
		}

		tok, err := svc.LookupToken(c.Context(), id)
		if err != nil {
			// For security, do not reveal whether the token was invalid or lookup failed.
			return fiber.ErrForbidden
		}

		c.Locals(localKeyToken, tok)
		return c.Next()
	}
}

func TokenFromContext(c fiber.Ctx) *token.Token[*SessionData] {
	raw := c.Locals(localKeyToken)
	tok, ok := raw.(*token.Token[*SessionData])
	if !ok {
		panic(fmt.Sprintf("auth.TokenFromContext: unexpected token type: %T", raw))
	}

	return tok
}
