package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"os"

	"github.com/SunnyRaj84348/do-notes/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RandToken() []byte {
	b := make([]byte, 64)

	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

type Credential struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	UserID   int
	Username string
	Password string
}

func main() {
	router := gin.Default()

	store := cookie.NewStore(RandToken())
	store.Options(sessions.Options{Secure: false})

	router.Use(sessions.Sessions("session_user", store))

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
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"message": "user already exist",
			})
		}
	})

	router.POST("/login", func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		val := session.Get("user")
		if val != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "user already logged in",
			})
		}

		cred := Credential{}

		err := ctx.BindJSON(&cred)
		if err != nil {
			return
		}

		user := User{}
		row := database.GetUser(db, cred.Username)

		err = row.Scan(&user.UserID, &user.Username, &user.Password)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password))
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		session.Set("user", user.Username)

		err = session.Save()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	})

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Could not start the http server: %v", err)
	}
}
