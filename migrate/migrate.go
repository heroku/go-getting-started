package main

import (
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/models"
)

func init() {
	// inits.LoadEnvVariables()
	inits.ConDB()
}

func main() {
	inits.DB.AutoMigrate(&models.User{})
	inits.DB.AutoMigrate(&models.TaskModel{})
	inits.DB.AutoMigrate(&models.Project{})
	inits.DB.AutoMigrate(&models.SubTaskModel{})
}

// func main() {
// 	inits.DB.AutoMigrate(&models.User{})
// 	inits.DB.AutoMigrate(&models.Project{})   // Ensure Project table is created first
// 	inits.DB.AutoMigrate(&models.TaskModel{}) // Then TaskModel
// 	// inits.DB.AutoMigrate(&models.SubTaskModel{}) // Finally SubTaskModel

// 	// Additional migration logic if needed
// }
