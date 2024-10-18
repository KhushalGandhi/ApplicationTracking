package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"net/http"
	"recruitment-system/models"
	"strings"
)

func GetUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the JWT token from the request header
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(http.StatusUnauthorized).SendString("No token provided")
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_secret_key"), nil
		})
		if err != nil || !token.Valid {
			return c.Status(http.StatusUnauthorized).SendString("Invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(http.StatusUnauthorized).SendString("Invalid token claims")
		}

		email := claims["email"].(string)

		// Retrieve the user from the database using the email from the JWT token
		var user models.User
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			return c.Status(http.StatusNotFound).SendString("User not found")
		}

		return c.Status(http.StatusOK).JSON(user)
	}
}

func GetAllUsers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error fetching users")
		}
		return c.Status(http.StatusOK).JSON(users)
	}
}
