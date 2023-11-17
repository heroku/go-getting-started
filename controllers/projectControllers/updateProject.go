package projectControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
)

func UpdateProject(c *gin.Context) {
	var modelProject models.Project
	var updateProject models.Project

	// Manually Bind Json the JSON data into the newProject struct
	if err := c.ShouldBindJSON(&modelProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the project ID is provided
	if modelProject.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	// Update a project record in the database
	if err := inits.DB.First(&updateProject, modelProject.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	// Update project fields
	if modelProject.ProjectName != "" {
		updateProject.ProjectName = modelProject.ProjectName
	}

	// Save the updated project
	inits.DB.Save(&updateProject)

	c.JSON(http.StatusOK, "update project successfully!")
}

func UpdateProjectByID(c *gin.Context) {
	var modelProject models.Project
	var updateProject models.Project

	// Get the project ID from the URL parameter
	projectID := c.Param("project_id")

	// Manually Bind Json the JSON data into the newProject struct
	if err := c.ShouldBindJSON(&modelProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the project ID is provided
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	// Update a project record in the database
	if err := inits.DB.Where("id = ?", projectID).First(&updateProject).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	// Update project fields
	if modelProject.ProjectName != "" {
		updateProject.ProjectName = modelProject.ProjectName
	}

	// Save the updated project
	inits.DB.Save(&updateProject)

	c.JSON(http.StatusOK, "update project successfully!")
}
