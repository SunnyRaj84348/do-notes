package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SunnyRaj84348/do-notes/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Credential struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	router := gin.Default()

	db, err := database.Connect(os.Getenv("MYSQL_STR"))
	if err != nil {
		log.Fatal(err)
	}

	router.POST("/signup", func(ctx *gin.Context) {
		cred := Credential{}

		err := ctx.BindJSON(&cred)
		if err != nil {
			return
		}

		hashPass, err := bcrypt.GenerateFromPassword([]byte(cred.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		err = database.InsertUser(db, cred.Username, string(hashPass))
		if err != nil {
			ctx.JSON(http.StatusConflict, gin.H{
				"message": "user already exist",
			})

			return
		}
	})

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Could not start the http server: %v", err)
	}
}
