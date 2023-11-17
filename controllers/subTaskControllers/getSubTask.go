package subTaskControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
)

func GetSubTask(c *gin.Context) {
	// Get the task ID from the URL parameter
	// projectID := c.Param("project_id")
	taskID := c.Param("task_id")

	var subTasks []models.SubTaskModel

	// Find all subTasks associated with the task
	if err := inits.DB.Where("task_sub_task_id = ?", taskID).Find(&subTasks).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subTasks not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sub_task_data": subTasks})
}

func GetSubTaskByID(c *gin.Context) {
	// Get the sub task ID from the URL parameter
	subTaskID := c.Param("sub_task_id")

	var subTask models.SubTaskModel

	// Find the sub task by ID
	if err := inits.DB.Where("id = ?", subTaskID).First(&subTask).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sub task not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sub_task_data": subTask})
}

// func GetSubTask(c *gin.Context) {
// 	// Get the task ID from the URL parameter
// 	projectID := c.Param("project_id")
// 	taskID := c.Param("task_id")

// 	var subTasks []models.SubTaskModel

// 	// Find all subTasks associated with the task
// 	if err := inits.DB.Where("project_id = ?", projectID).Where("task_id = ?", taskID).Find(&subTasks).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "subTasks not found!"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"sub_task_data": subTasks})
// }

// func GetSubTaskByID(c *gin.Context) {
// 	// Get the task ID from the URL parameter
// 	subTaskID := c.Param("sub_task_id")
// 	taskID := c.Param("task_id")

// 	var subTask models.SubTaskModel

// 	// Find the task by ID
// 	if err := inits.DB.Where("task_id = ?", taskID).Where("id = ?", subTaskID).First(&subTask).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Sub task not found!"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"sub_task_data": subTask})
// }
