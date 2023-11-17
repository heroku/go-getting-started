package taskControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
)

func UpdateTask(c *gin.Context) {
	var modelTask models.TaskModel
	var updateTask models.TaskModel

	// Manually Bind Json the JSON data into the newTask struct
	if err := c.ShouldBindJSON(&modelTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the task ID is provided
	if modelTask.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	// Update a new task record in the database
	if err := inits.DB.First(&updateTask, modelTask.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Update task fields
	if modelTask.TaskName != "" {
		updateTask.TaskName = modelTask.TaskName
	}
	if modelTask.Status != "" {
		updateTask.Status = modelTask.Status
	}
	if modelTask.Description != "" {
		updateTask.Description = modelTask.Description
	}
	if modelTask.UserTaskID != 0 {
		updateTask.UserTaskID = modelTask.UserTaskID
	}

	// Save the updated task
	inits.DB.Save(&updateTask)

	c.JSON(http.StatusOK, "update task successfully!")
}
func UpdateTaskByID(c *gin.Context) {
	var modelTask models.TaskModel
	var updateTask models.TaskModel

	// Get the project ID from the URL parameter
	// projectID := c.Param("project_id")

	// Get the task ID from the URL parameter
	taskID := c.Param("task_id")

	// Manually Bind Json the JSON data into the newTask struct
	if err := c.ShouldBindJSON(&modelTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the task ID is provided
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	// Update a task record in the database
	if err := inits.DB.Where("id = ?", taskID).First(&updateTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Update task fields
	if modelTask.TaskName != "" {
		updateTask.TaskName = modelTask.TaskName
	}
	if modelTask.Status != "" {
		updateTask.Status = modelTask.Status
	}
	if modelTask.Description != "" {
		updateTask.Description = modelTask.Description
	}
	if modelTask.UserTaskID != 0 {
		updateTask.UserTaskID = modelTask.UserTaskID
	}

	// Save the updated task
	inits.DB.Save(&updateTask)

	c.JSON(http.StatusOK, "update task successfully!")
}
