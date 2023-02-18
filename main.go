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

// Generate 64 secure random numbers
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

type Note struct {
	NoteTitle string `json:"noteTitle" binding:"required"`
	NoteBody  string `json:"noteBody" binding:"required"`
}

func main() {
	router := gin.Default()

	// Create new cookie store with secure auth
	store := cookie.NewStore(RandToken())
	store.Options(sessions.Options{Secure: false})

	router.Use(sessions.Sessions("session_user", store))

	// Init database connection
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

		// Hash the given password using bcrypt
		hashPass, err := bcrypt.GenerateFromPassword([]byte(cred.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Finally, insert user cred to database
		err = database.InsertUser(db, cred.Username, string(hashPass))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"message": "user already exist",
			})
			return
		}
	})

	router.POST("/login", func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		// Check for existing session
		val := session.Get("user")
		if val != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "user already logged in",
			})
			return
		}

		cred := Credential{}

		err := ctx.BindJSON(&cred)
		if err != nil {
			return
		}

		user := User{}
		row := database.GetUser(db, cred.Username)

		// Match username with database
		err = row.Scan(&user.UserID, &user.Username, &user.Password)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Check for invalid password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password))
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Add user to the session
		session.Set("user", user.Username)

		err = session.Save()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	})

	router.POST("/logout", func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		// Check if session doesn't exist
		val := session.Get("user")
		if val == nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		session.Delete("user")

		err := session.Save()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	})

	router.POST("/insert-note", func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		// Check if session doesn't exist
		username := session.Get("user")
		if username == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		note := Note{}

		err := ctx.BindJSON(&note)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Insert note into database
		err = database.InsertNotes(db, username.(string), note.NoteTitle, note.NoteBody)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	})

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Could not start the http server: %v", err)
	}
}
