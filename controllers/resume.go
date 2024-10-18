package controllers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"recruitment-system/models"
	"strings"
)

const resumeParserAPI = "https://api.apilayer.com/resume_parser/upload"
const apiKey = "0bWeisRWolJ3UdX3MXMSMWpYfPIpQfS"

func UploadResume(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if the user is authenticated
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
		}
		// Handle file upload
		file, err := c.FormFile("resume")
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("File not found")
		}
		fileBytes, err := file.Open()
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error reading file")
		}
		defer fileBytes.Close()

		// Send file to third-party API
		req, err := http.NewRequest("POST", resumeParserAPI, fileBytes)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error creating request")
		}
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("apikey", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error calling resume parser API")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error reading API response")
		}

		// Parse the API response
		var resumeData map[string]interface{}
		if err := json.Unmarshal(body, &resumeData); err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error parsing API response")
		}

		//var user models.User
		// Implement logic to get the user from the database
		// For now, assuming user ID is hardcoded
		userID := uint(1)

		profile := models.Profile{
			UserID:     userID,
			ResumeFile: file.Filename,
			Skills:     getStringFromData(resumeData, "skills"),
			Education:  getStringFromData(resumeData, "education"),
			Experience: getStringFromData(resumeData, "experience"),
		}

		if err := db.Create(&profile).Error; err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error saving profile")
		}

		return c.Status(http.StatusOK).SendString("Resume uploaded successfully")
	}
}

func getStringFromData(data map[string]interface{}, key string) string {
	if value, exists := data[key]; exists {
		switch v := value.(type) {
		case []interface{}:
			var result []string
			for _, item := range v {
				if str, ok := item.(string); ok {
					result = append(result, str)
					// Handle nested objects if needed
				}
			}
			return strings.Join(result, ", ")
		case string:
			return v
		}
	}
	return ""
}
