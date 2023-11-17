package userControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
)

func AddProject(c *gin.Context) {
	var newProject models.Project
	userID := c.Param("user_id")
	var existUser models.User

	// Manually Bind Json the JSON data into the newProject struct
	if err := c.ShouldBindJSON(&newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if project_name is empty
	if newProject.ProjectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_name is required"})
		return
	}

	// Get the user ID from the request context (assuming it's stored there)
	inits.DB.Where("id = ?", userID).First(&existUser)

	// userID, exists := c.Get("user_id")
	if existUser.ID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})
		return
	}

	// Create a project in the database
	if err := inits.DB.Create(&newProject).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Project"})
		return
	} else {
		// Associate the user with the project
		inits.DB.Model(&newProject).Association("Users").Append(&existUser)
		// inits.DB.Model(&existUser).Association("Projects").Append(&newProject)

		c.JSON(http.StatusOK, gin.H{
			"message": "created new Project!",
			"add_by":  existUser.Username,
			"id":      existUser.ID,
		})
	}
}
