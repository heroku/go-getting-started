package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/inits"
	"github.com/heroku/go-getting-started/routes"
	"github.com/joho/godotenv"
)

func init() {
	inits.LoadEnvVariables()
	inits.ConDB()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	//call gin function another example => gin.New()
	r := gin.Default()
	r.Use(cors.Default())
	routes.Routes(r)

	//*gin.Context is default
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	r.Run() // listen and serve on 0.0.0.0:8080
}
