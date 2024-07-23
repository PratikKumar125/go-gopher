package guards

import (
	"first/models"
	"first/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JwtAuthGuard(ctx *fiber.Ctx) error {
	var defaultTokenLookup = "header:" + fiber.HeaderAuthorization
	parts := strings.Split(strings.TrimSpace(defaultTokenLookup), ":")
	tokenString, err := utils.JwtFromHeader(parts[1], ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error parsing jwt")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return utils.GetJWTSecret(), nil
	})

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid JWT token",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIdStr := claims["_id"].(string)
		userObjectID, _ := primitive.ObjectIDFromHex(userIdStr)

		user := models.User{
			ID:    userObjectID,
			Name:  claims["name"].(string),
			Email: claims["email"].(string),
		}
		ctx.Locals("user", user)
	}
	return ctx.Next()
}