package taskControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
)

func GetTask(c *gin.Context) {
	// Get the project ID from the URL parameter
	projectID := c.Param("project_id")

	var tasks []models.TaskModel

	// Find all tasks associated with the project
	if err := inits.DB.Where("project_task_id = ?", projectID).Preload("SubTask").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tasks not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task_data": tasks})
}

func GetTaskByID(c *gin.Context) {
	// Get the task ID from the URL parameter
	taskID := c.Param("task_id")
	projectID := c.Param("project_id")

	var task models.TaskModel

	// Find the task by ID
	if err := inits.DB.Where("project_task_id = ?", projectID).Where("id = ?", taskID).Preload("SubTask").First(&task).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task_data": task})
}

// func GetTask(c *gin.Context) {
// 	// Get the project ID from the URL parameter
// 	projectID := c.Param("project_id")

// 	var tasks []models.TaskModel

// 	// Find all tasks associated with the project
// 	if err := inits.DB.Where("project_id = ?", projectID).Find(&tasks).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Tasks not found!"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"task_data": tasks})
// }

// func GetTaskByID(c *gin.Context) {
// 	// Get the task ID from the URL parameter
// 	taskID := c.Param("task_id")
// 	projectID := c.Param("project_id")

// 	var task models.TaskModel

// 	// Find the task by ID
// 	if err := inits.DB.Where("project_id = ?", projectID).Where("id = ?", taskID).First(&task).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Task not found!"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"task_data": task})
// }
