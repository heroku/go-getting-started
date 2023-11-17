package subTaskControllers

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
)

func DeleteSubTaskByID(c *gin.Context) {
	var modelSubTask models.SubTaskModel
	// taskID := c.Param("task_id")

	subTaskID := c.Param("sub_task_id")

	if err := inits.DB.Where("id = ?", subTaskID).First(&modelSubTask).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maybe the sub task has already been deleted!"})
		return
	}

	if err := inits.DB.Delete(&modelSubTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, "Deleted sub task with ID: "+subTaskID)
}
