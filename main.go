package main

import (
	"log"
	"os"

	"github.com/SunnyRaj84348/do-notes/controllers"
	"github.com/SunnyRaj84348/do-notes/middlewares"
	"github.com/SunnyRaj84348/do-notes/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()

	// Use middleware
	router.Use(middlewares.Sessions())
	router.Use(middlewares.Cors())

	// Load .env file vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Init database connection
	err = models.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	// Init database tables
	err = models.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// Load html template files
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")

	auth := router.Group("/", middlewares.Auth)
	{
		auth.GET("/", controllers.GetNotes)
		auth.POST("/logout", controllers.Logout)
		auth.POST("/insert-note", controllers.InsertNote)
		auth.PUT("/update-note/:id", controllers.UpdateNote)
		auth.DELETE("/delete-note/:id", controllers.DeleteNote)
	}

	router.GET("/login", controllers.Login)
	router.POST("/session", controllers.Session)
	router.POST("/signup", controllers.Signup)

	err = router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Could not start the http server: %v", err)
	}
}
