package main

import (
	"crypto/rand"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/SunnyRaj84348/do-notes/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	NoteBody  string `json:"noteBody"`
}

type Notes struct {
	NoteID    int    `json:"noteID"`
	NoteTitle string `json:"noteTitle"`
	NoteBody  string `json:"noteBody"`
}

func main() {
	router := gin.Default()

	// Create new cookie store with secure auth
	store := cookie.NewStore(RandToken())
	store.Options(sessions.Options{Secure: false})

	router.Use(sessions.Sessions("session_user", store))
	router.Use(cors.Default())

	// Load .env file vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Init database connection
	db, err := database.Connect(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	err = database.Init(db)
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
			if err == sql.ErrNoRows {
				ctx.AbortWithStatus(http.StatusUnauthorized)
			} else {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}

			return
		}

		// Check for invalid password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password))
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Add user to the session
		session.Set("user", user.UserID)

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

		// Delete session_user cookie
		session.Clear()
		session.Options(sessions.Options{MaxAge: -1})

		err := session.Save()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	})

	router.POST("/insert-note", func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		// Check if session doesn't exist
		userid := session.Get("user")
		if userid == nil {
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
		noteID, err := database.InsertNotes(db, userid.(int), note.NoteTitle, note.NoteBody)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		notes := Notes{noteID, note.NoteTitle, note.NoteBody}

		// Write inserted note to response body
		ctx.JSON(http.StatusOK, notes)
	})

	router.GET("/get-notes", func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		// Check if session doesn't exist
		userid := session.Get("user")
		if userid == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Retrieve notes from database
		rows, err := database.GetNotes(db, userid.(int))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		defer rows.Close()

		notes := []Notes{}

		for rows.Next() {
			note := Notes{}

			err := rows.Scan(&note.NoteID, &note.NoteTitle, &note.NoteBody)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			notes = append(notes, note)
		}

		// Write notes to response body
		ctx.JSON(http.StatusOK, notes)
	})

	router.PUT("/update-note/:id", func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		// Check if session doesn't exist
		userid := session.Get("user")
		if userid == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		note := Note{}

		err := ctx.BindJSON(&note)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		val := ctx.Param("id")

		// Convert string to int
		id, err := strconv.Atoi(val)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Update specified note
		err = database.UpdateNotes(db, userid.(int), id, note.NoteTitle, note.NoteBody)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.AbortWithStatus(http.StatusForbidden)
			} else {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}

			return
		}

		notes := Notes{id, note.NoteTitle, note.NoteBody}

		// Write updated note to response body
		ctx.JSON(http.StatusOK, notes)
	})

	router.DELETE("/delete-note/:id", func(ctx *gin.Context) {
		session := sessions.Default(ctx)

		// Check if session doesn't exist
		userid := session.Get("user")
		if userid == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		val := ctx.Param("id")

		// Convert string to int
		id, err := strconv.Atoi(val)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Delete specified note from database
		err = database.DeleteNotes(db, userid.(int), id)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.AbortWithStatus(http.StatusForbidden)
			} else {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}

			return
		}
	})

	err = router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Could not start the http server: %v", err)
	}
}
