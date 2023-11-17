package subTaskControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
)

func AddSubTask(c *gin.Context) {
	var newSubTask models.SubTaskModel

	// Manually Bind Json the JSON data into the newSubTask struct
	if err := c.ShouldBindJSON(&newSubTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if task_name is empty
	if newSubTask.TaskName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sub_task_name is required"})
		return
	}

	// projectID := c.Param("project_id")
	taskID := c.Param("task_id") // Get the task ID from the URL parameter

	var task models.TaskModel

	// Check if the task with the given IDs exists
	if err := inits.DB.Where("id = ?", taskID).First(&task).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task not found or not in project"})
		return
	}

	// Set the taskID for the new task
	newSubTask.TaskSubTaskID = task.ID

	// Create a new task record in the database
	if err := inits.DB.Create(&newSubTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sub task", "err": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "created new sub task!", "sub_task_data": newSubTask})
	}
}

// func AddSubTask(c *gin.Context) {
// 	var newSubTask models.SubTaskModel

// 	// Manually Bind Json the JSON data into the newSubTask struct
// 	if err := c.ShouldBindJSON(&newSubTask); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Check if task_name is empty
// 	if newSubTask.TaskName == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "sub_task_name is required"})
// 		return
// 	}

// 	projectID := c.Param("project_id")
// 	taskID := c.Param("task_id") // Get the task ID from the URL parameter

// 	var task models.TaskModel

// 	// Check
// 	if err := inits.DB.Where("project_id = ?", projectID).Where("id = ?", taskID).First(&task).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "not in project"})
// 		return
// 	}

// 	// Check if the task with the given ID exists
// 	if err := inits.DB.Where("id = ?", taskID).First(&task).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "task not found!"})
// 		return
// 	}

// 	// Set the taskID for the new task
// 	newSubTask.TaskSubTaskID = task.ID

// 	// Create a new task record in the database
// 	if err := inits.DB.Create(&newSubTask).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sub task", "err": err.Error()})
// 		return
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{"message": "created new sub task!"})
// 	}
// }
