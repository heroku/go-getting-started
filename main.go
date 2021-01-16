package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Joojo7/go-getting-started/router"
	"github.com/gorilla/mux"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
)

func main() {
	// port := os.Getenv("PORT")
	// fmt.Print(port)

	// if port == "" {
	// 	port = "8000"
	// }

	// router := gin.New()
	// router.Use(gin.Logger())
	// router.LoadHTMLGlob("templates/*.tmpl.html")
	// router.Static("/static", "static")

	// router.GET("/tut", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.tmpl.html", nil)
	// })

	myEnv, err1 := godotenv.Read()
	if err1 != nil {
		log.Fatal(err1)
	}
	port := myEnv["PORT"]

	myRouter := mux.NewRouter().StrictSlash(true)

	// ROuter files
	router.Routes(myRouter)
	router.FoodRoutes(myRouter)
	router.OrderItemRoutes(myRouter)
	router.TableRoutes(myRouter)
	router.InvoiceRoutes(myRouter)
	router.OrderRoutes(myRouter)

	fmt.Printf("listening on %v \n", port)
	error1 := http.ListenAndServe(port, myRouter)
	if error1 != nil {
		panic(error1)
	}
}
