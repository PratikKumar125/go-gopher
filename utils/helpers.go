package utils

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrJWTMissingOrMalformed = errors.New("missing or malformed JWT")
	DefaultAuthScheme = "Bearer"
)

func JwtFromHeader(header string, c *fiber.Ctx) (string, error) {
	auth := c.Get(header)
	l := len(DefaultAuthScheme)
	if len(auth) > l+1 && strings.EqualFold(auth[:l], DefaultAuthScheme) {
		return strings.TrimSpace(auth[l:]), nil
	}
	return "", ErrJWTMissingOrMalformed
}

func SignJwtToken(claims *jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"));
	return t, err
}