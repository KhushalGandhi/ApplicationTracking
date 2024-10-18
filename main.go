package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"recruitment-system/controllers"
	"recruitment-system/middleware"
	"recruitment-system/models"
)

func main() {
	app := fiber.New()
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto migrate models
	db.AutoMigrate(&models.User{}, &models.Profile{}, &models.Job{}, &models.Application{})

	// Routes
	app.Post("/signup", controllers.SignUp(db))
	app.Post("/login", controllers.Login(db))
	app.Post("/uploadResume", middleware.AuthMiddleware, controllers.UploadResume(db))
	app.Post("/admin/job", middleware.AuthMiddleware, controllers.CreateJob(db))
	app.Get("/admin/job/:job_id", middleware.AuthMiddleware, controllers.GetJob(db))
	app.Get("/admin/applicants", middleware.AuthMiddleware, controllers.GetApplicants(db))
	app.Get("/admin/applicant/:applicant_id", middleware.AuthMiddleware, controllers.GetApplicant(db))
	app.Get("/jobs", controllers.GetJobs(db))
	app.Post("/jobs/apply", middleware.AuthMiddleware, controllers.ApplyJob(db)) // Implement ApplyJob controller if needed

	app.Listen(":3000")
}
