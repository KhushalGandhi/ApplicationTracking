package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"net/http"
	"recruitment-system/models"
	"strings"
)

func CreateJob(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var job models.Job
		if err := c.BodyParser(&job); err != nil {
			log.Error(err)

			return c.Status(http.StatusBadRequest).SendString("Invalid request")
		}
		// Check if the user is an admin
		token := c.Get("Authorization")
		if token == "" {
			log.Error("unauthorized")

			return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
		}
		// Implement authentication and authorization here

		if err := db.Create(&job).Error; err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error creating job")
		}
		return c.Status(http.StatusOK).JSON(job)
	}
}

func GetJob(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobID := c.Params("job_id")
		var job models.Job
		if err := db.Preload("PostedBy").Where("id = ?", jobID).First(&job).Error; err != nil {
			log.Error(err)

			return c.Status(http.StatusNotFound).JSON("Job not found")
		}
		return c.Status(http.StatusOK).JSON(job)
	}
}

func GetJobs(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var jobs []models.Job

		// Find all jobs and preload the PostedBy relationship
		if err := db.Preload("PostedBy").Find(&jobs).Error; err != nil {
			log.Error(err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}

		if len(jobs) == 0 {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "No jobs found",
			})
		}

		return c.Status(http.StatusOK).JSON(jobs)
	}
}

func GetApplicants(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			log.Error(err)

			return c.Status(http.StatusInternalServerError).SendString("Error fetching applicants")
		}
		return c.Status(http.StatusOK).JSON(users)
	}
}

func GetApplicant(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		applicantID := c.Params("applicant_id")
		var profile models.Profile
		if err := db.Where("user_id = ?", applicantID).First(&profile).Error; err != nil {
			log.Error(err)

			return c.Status(http.StatusNotFound).SendString("Applicant not found")
		}
		return c.Status(http.StatusOK).JSON(profile)
	}
}

func ApplyJob(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var application struct {
			JobID uint `json:"job_id"`
		}
		if err := c.BodyParser(&application); err != nil {
			log.Error(err)

			return c.Status(http.StatusBadRequest).SendString("Invalid request")
		}

		// Get the JWT token from the request header
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			log.Error("no token provided")

			return c.Status(http.StatusUnauthorized).SendString("No token provided")
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			return []byte("your_secret_key"), nil
		})
		if err != nil || !token.Valid {
			log.Error(err)

			return c.Status(http.StatusUnauthorized).SendString("Invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Error(err)

			return c.Status(http.StatusUnauthorized).SendString("Invalid token claims")
		}

		email := claims["email"].(string)

		// Retrieve the user from the database using the email from the JWT token
		var user models.User
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			log.Error(err)

			return c.Status(http.StatusNotFound).SendString("User not found")
		}

		// Check if the job exists
		var job models.Job
		if err := db.Where("id = ?", application.JobID).First(&job).Error; err != nil {
			log.Error(err)

			return c.Status(http.StatusNotFound).SendString("Job not found")
		}

		// Check if the user has already applied for this job
		var existingApplication models.Application
		if err := db.Where("user_id = ? AND job_id = ?", user.ID, job.ID).First(&existingApplication).Error; err == nil {
			log.Error(err)

			return c.Status(http.StatusConflict).SendString("User has already applied for this job")
		}

		// Create a new application
		newApplication := models.Application{
			UserID: user.ID,
			JobID:  job.ID,
		}
		if err := db.Create(&newApplication).Error; err != nil {
			log.Error(err)

			return c.Status(http.StatusInternalServerError).SendString("Error applying for job")
		}

		// Update the total number of applications for the job
		job.TotalApplications++
		if err := db.Save(&job).Error; err != nil {
			log.Error(err)

			return c.Status(http.StatusInternalServerError).SendString("Error updating job applications count")
		}

		return c.Status(http.StatusOK).JSON(newApplication)
	}
}
