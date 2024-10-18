package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		log.Error("No token provided")
		return c.Status(http.StatusUnauthorized).SendString("No token provided")
	}
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})
	if err != nil || !token.Valid {
		log.Error("invalid token")
		return c.Status(http.StatusUnauthorized).SendString("Invalid token")
	}
	return c.Next()
}
