package main

import (
	"log"
	"os"

	"github.com/SunnyRaj84348/do-notes/database"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	db, err := database.Connect(os.Getenv("MYSQL_STR"))
	if err != nil {
		log.Fatal(err)
	}

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Could not start the http server: %v", err)
	}
}
