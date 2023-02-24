package controllers

import (
	"database/sql"
	"net/http"

	"github.com/SunnyRaj84348/do-notes/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(ctx *gin.Context) {
	cred := models.Credential{}

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

	// Insert user cred to database
	err = models.InsertUser(cred.Username, string(hashPass))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "user already exist",
		})
		return
	}
}

func Login(ctx *gin.Context) {
	cred := models.Credential{}

	err := ctx.BindJSON(&cred)
	if err != nil {
		return
	}

	user := models.User{}
	row := models.GetUser(cred.Username)

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

	session := sessions.Default(ctx)

	// Add user to the session
	session.Set("user", user.UserID)

	err = session.Save()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	// Delete session_user cookie
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})

	err := session.Save()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
