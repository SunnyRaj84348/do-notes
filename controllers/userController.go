package controllers

import (
	"net/http"

	"github.com/SunnyRaj84348/do-notes/mailjet"
	"github.com/SunnyRaj84348/do-notes/models"
	"github.com/SunnyRaj84348/do-notes/utilities"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(ctx *gin.Context) {
	user := models.User{}

	err := ctx.BindJSON(&user)
	if err != nil {
		return
	}

	err = utilities.ValidateEmail(user.Email)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Hash the given password using bcrypt
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user.Password = string(hashPass)

	// Insert user cred to database
	err = models.InsertUser(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": "user already exist",
		})
		return
	}

	code, err := utilities.GenVerificationCode()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = mailjet.SendVerification(user.Email, code)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = models.InsertEmailAuth(user.Email, code)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func Login(ctx *gin.Context) {
	cred := models.Credential{}

	err := ctx.BindJSON(&cred)
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

	// Check if email is not verified
	if !user.Verified {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "email has not been verified",
		})
		return
	}

	session := sessions.Default(ctx)

	// Add user to the session
	session.Set("user", user.UserID)

	err = session.Save()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
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
	}
}

func VerifyEmail(ctx *gin.Context) {
	emailAuth := models.EmailAuth{}

	err := ctx.BindJSON(&emailAuth)
	if err != nil {
		return
	}

	err = utilities.ValidateEmail(emailAuth.Email)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = models.GetEmailAuth(emailAuth)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = models.DeleteEmailAuth(emailAuth)
	if err != nil {
		ctx.Copy().AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = models.SetUserVerified(emailAuth.Email)
	if err != nil {
		ctx.Copy().AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
