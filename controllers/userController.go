package controllers

import (
	"net/http"

	"github.com/SunnyRaj84348/do-notes/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	}
}

func Session(ctx *gin.Context) {
	cred := models.Credential{}

	err := ctx.Bind(&cred)
	if err != nil {
		return
	}

	user, err := models.GetUser(cred.Username)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
	}

	// Redirect to homepage
	ctx.Redirect(http.StatusFound, "/")
}

func Login(ctx *gin.Context) {
	session := sessions.Default(ctx)

	userid := session.Get("user")
	if userid != nil {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	ctx.HTML(http.StatusOK, "login.html", nil)
}

func Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	// Delete session_user cookie
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})

	err := session.Save()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}
