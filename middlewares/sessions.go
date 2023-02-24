package middlewares

import (
	"github.com/SunnyRaj84348/do-notes/utilities"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Sessions() gin.HandlerFunc {
	// Create new cookie store with secure auth
	store := cookie.NewStore(utilities.RandToken())
	store.Options(sessions.Options{
		Secure: false,
		MaxAge: 60 * 60 * 24 * 30,
	})

	return sessions.Sessions("session_user", store)
}
