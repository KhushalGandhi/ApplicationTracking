package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"net/http"
	"recruitment-system/models"
	"time"
)

var jwtSecret = []byte("your_secret_key")

func SignUp(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User
		if err := c.BodyParser(&user); err != nil {
			log.Error(err)
			return c.Status(http.StatusBadRequest).SendString("Invalid request")
		}
		////hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		////if err != nil {
		////	log.Error(err)
		////	return c.Status(http.StatusInternalServerError).SendString("Error hashing password")
		////}
		//user.PasswordHash = string(hash)
		if err := db.Create(&user).Error; err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString("Error creating user")
		}
		return c.Status(http.StatusOK).JSON(user)
	}
}

func Login(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&requestBody); err != nil {
			log.Error(err)
			return c.Status(http.StatusBadRequest).SendString("Invalid request")
		}
		var user models.User
		if err := db.Where("email = ?", requestBody.Email).First(&user).Error; err != nil {
			log.Error(err)
			return c.Status(http.StatusUnauthorized).SendString("User not found")
		}
		//if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(requestBody.Password)); err != nil {
		//	log.Error(err)
		//	return c.Status(http.StatusUnauthorized).SendString("Invalid credentials")
		//}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).SendString("Error generating token")
		}
		return c.Status(http.StatusOK).JSON(fiber.Map{"token": tokenString})
	}
}
