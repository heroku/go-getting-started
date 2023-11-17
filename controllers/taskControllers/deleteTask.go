package taskControllers

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
)

func DeleteTaskByID(c *gin.Context) {
	var modelTask models.TaskModel
	// projectID := c.Param("project_id")
	taskID := c.Param("task_id")

	if err := inits.DB.Where("id = ?", taskID).First(&modelTask).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Maybe the task has already been deleted!"})
		return
	}

	if err := inits.DB.Delete(&modelTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, "Deleted task with ID: "+taskID)
}
