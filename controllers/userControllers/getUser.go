package userControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
)

func GetUser(c *gin.Context) {
	// []models call array (all elements)
	var users []models.User

	if err := inits.DB.Preload("UserTask").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	// Fetch associated tasks and users for each project
	for i := range users {
		var projects []models.Project

		inits.DB.Model(&users[i]).Association("Projects").Find(&projects)
		users[i].Projects = projects
	}
	c.JSON(http.StatusOK, gin.H{"user_data": users})
}

func GetUserByID(c *gin.Context) {
	// models call one element
	var user models.User

	if err := inits.DB.Where("id = ?", c.Param("user_id")).Preload("UserTask").First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Fetch associated tasks and users for each project
	var projects []models.Project

	inits.DB.Model(&user).Association("Projects").Find(&projects)
	user.Projects = projects

	c.JSON(http.StatusOK, gin.H{"user_data": user})
}
