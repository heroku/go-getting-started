package projectControllers

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
)

func GetProject(c *gin.Context) {
	var projects []models.Project

	// Find all projects
	if err := inits.DB.Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects"})
		return
	}

	// Fetch associated tasks and users for each project
	for i := range projects {
		var tasks []models.TaskModel
		inits.DB.Where("project_task_id = ?", projects[i].ID).Find(&tasks)
		projects[i].Task = tasks

		var users []models.User
		inits.DB.Model(&projects[i]).Association("Users").Find(&users)
		projects[i].Users = users
	}

	c.JSON(http.StatusOK, gin.H{"projects_data": projects})
}

// func GetProjectByID(c *gin.Context) {
// 	// get the posts
// 	// models call one element
// 	var project models.TaskModel

// 	if err := inits.DB.Where("id = ?", c.Param("id")).First(&project).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
// 		return
// 	}

//		c.JSON(http.StatusOK, gin.H{"project_data": project})
//	}
func GetProjectByID(c *gin.Context) {
	// Get the project ID from the URL parameter
	projectID := c.Param("project_id")

	// Find the project by ID
	var project models.Project
	if err := inits.DB.Where("id = ?", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project not found!"})
		return
	}

	// Fetch associated tasks for the project
	var tasks []models.TaskModel
	inits.DB.Where("project_id = ?", project.ID).Find(&tasks)
	project.Task = tasks

	var users []models.User
	inits.DB.Model(&project).Association("Users").Find(&users)
	project.Users = users

	c.JSON(http.StatusOK, gin.H{"project_data": project})
}
