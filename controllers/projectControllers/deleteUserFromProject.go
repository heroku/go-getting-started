package projectControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
)

func DeleteUserFromProject(c *gin.Context) {
	userID := c.Param("user_id")
	projectID := c.Param("project_id")

	var project models.Project
	var user models.User

	// Check if project with projectID exists
	if err := inits.DB.Where("id = ?", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project not found"})
		return
	}

	// Check if user with userID exists
	if err := inits.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Remove the user from the project
	inits.DB.Model(&project).Association("Users").Delete(&user)

	// Fetch associated tasks for the project
	var tasks []models.TaskModel
	inits.DB.Where("project_id = ?", project.ID).Find(&tasks)
	project.Task = tasks

	c.JSON(http.StatusOK, gin.H{"project_data": project})
}
