package userControllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check duplicate User
	var duplicateUser models.User
	inits.DB.Where("username = ?", user.Username).First(&duplicateUser)
	if duplicateUser.ID > 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "Username already exists"})
		return
	}
	// Hash (encrypt) the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)
	// Save the user to the database
	inits.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"message":  "user registered successfully",
		"username": user.Username,
		"ID":       user.ID,
	})
}
