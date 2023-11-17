package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/controllers/projectControllers"
	"github.com/heroku/go-getting-started/controllers/subTaskControllers"
	"github.com/heroku/go-getting-started/controllers/taskControllers"
	"github.com/heroku/go-getting-started/controllers/userControllers"
)

func Routes(router *gin.Engine) {
	// user auth routes
	router.POST("/register", userControllers.Register)
	router.POST("/login", userControllers.Login)

	// Apply the middleware to protected routes
	// router.Use(middlewares.AuthMiddleware())
	// router.Use(cors.Default())

	router.GET("/user", userControllers.GetUser)
	router.GET("/user/:user_id", userControllers.GetUserByID)
	router.POST("/user/:user_id/project", userControllers.AddProject)

	// project routes
	router.GET("/user/:user_id/project", projectControllers.GetProject)
	router.GET("/user/:user_id/project/:project_id", projectControllers.GetProjectByID)
	router.PUT("/user/:user_id/project", projectControllers.UpdateProject)
	router.PUT("/user/:user_id/project/:project_id", projectControllers.UpdateProjectByID)
	router.DELETE("/user/:user_id/project/:project_id", projectControllers.DeleteProjectByID)
	router.POST("/user/:user_id/project/:project_id/user/:invite_id", projectControllers.AddUserToProject)
	router.DELETE("/project/:project_id/user/:user_id", projectControllers.DeleteUserFromProject)

	// task routes
	router.GET("/user/:user_id/project/:project_id/task", taskControllers.GetTask)
	router.GET("/user/:user_id/project/:project_id/task/:task_id", taskControllers.GetTaskByID)
	router.POST("/user/:user_id/project/:project_id/task", taskControllers.AddTask)
	router.PUT("/user/:user_id/project/:project_id/task", taskControllers.UpdateTask)
	router.PUT("/user/:user_id/project/:project_id/task/:task_id", taskControllers.UpdateTaskByID)
	router.DELETE("/user/:user_id/project/:project_id/task/:task_id", taskControllers.DeleteTaskByID)

	// sub-task routes
	router.GET("/user/:user_id/project/:project_id/task/:task_id/subtask", subTaskControllers.GetSubTask)
	router.GET("/user/:user_id/project/:project_id/task/:task_id/:sub_task_id", subTaskControllers.GetSubTaskByID)
	router.POST("/user/:user_id/project/:project_id/task/:task_id/subtask", subTaskControllers.AddSubTask)
	router.PUT("/user/:user_id/project/:project_id/task/:task_id/subtask", subTaskControllers.UpdateSubTask)
	router.PUT("/user/:user_id/project/:project_id/task/:task_id/:sub_task_id", subTaskControllers.UpdateSubTaskByID)
	router.DELETE("/user/:user_id/project/:project_id/task/:task_id/:sub_task_id", subTaskControllers.DeleteSubTaskByID)

}
