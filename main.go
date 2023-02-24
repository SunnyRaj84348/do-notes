package main

import (
	"log"
	"os"

	"github.com/SunnyRaj84348/do-notes/controllers"
	"github.com/SunnyRaj84348/do-notes/initializers"
	"github.com/SunnyRaj84348/do-notes/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Use middleware
	router.Use(middlewares.Sessions())
	router.Use(middlewares.Cors())

	// Load .env file vars
	err := initializers.LoadEnvVars()
	if err != nil {
		log.Fatal(err)
	}

	// Init database connection
	err = initializers.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	// Init database tables
	err = initializers.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	auth := router.Group("/", middlewares.Auth)
	{
		auth.POST("/logout", controllers.Logout)
		auth.POST("/insert-note", controllers.InsertNote)
		auth.GET("/get-notes", controllers.GetNotes)
		auth.PUT("/update-note/:id", controllers.UpdateNote)
		auth.DELETE("/delete-note/:id", controllers.DeleteNote)
	}

	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)

	err = router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Could not start the http server: %v", err)
	}
}
