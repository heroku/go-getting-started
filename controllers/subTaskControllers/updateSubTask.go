package subTaskControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
)

func UpdateSubTask(c *gin.Context) {
	var modelSubTask models.SubTaskModel
	var updateSubTask models.SubTaskModel

	// Manually Bind Json the JSON data into the newSubTask struct
	if err := c.ShouldBindJSON(&modelSubTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the sub task ID is provided
	if modelSubTask.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sub task ID is required"})
		return
	}

	// Update a new task record in the database
	if err := inits.DB.First(&updateSubTask, modelSubTask.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Update task fields
	if modelSubTask.TaskName != "" {
		updateSubTask.TaskName = modelSubTask.TaskName
	}
	if modelSubTask.Status != "" {
		updateSubTask.Status = modelSubTask.Status
	}
	if modelSubTask.Description != "" {
		updateSubTask.Description = modelSubTask.Description
	}

	// Save the updated sub task
	inits.DB.Save(&updateSubTask)

	c.JSON(http.StatusOK, "update sub task successfully!")
}
func UpdateSubTaskByID(c *gin.Context) {
	var modelSubTask models.SubTaskModel
	var updateSubTask models.SubTaskModel

	// Get the task ID from the URL parameter
	taskID := c.Param("task_id")

	// Get the sub task ID from the URL parameter
	subTaskID := c.Param("sub_task_id")

	// Manually Bind Json the JSON data into the newTask struct
	if err := c.ShouldBindJSON(&modelSubTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the task ID is provided
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sub task ID is required"})
		return
	}

	// Update a task record in the database
	if err := inits.DB.Where("id = ?", subTaskID).First(&updateSubTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Update task fields
	if modelSubTask.TaskName != "" {
		updateSubTask.TaskName = modelSubTask.TaskName
	}
	if modelSubTask.Status != "" {
		updateSubTask.Status = modelSubTask.Status
	}
	if modelSubTask.Description != "" {
		updateSubTask.Description = modelSubTask.Description
	}

	// Save the updated task
	inits.DB.Save(&updateSubTask)

	c.JSON(http.StatusOK, "update sub task successfully!")
}
