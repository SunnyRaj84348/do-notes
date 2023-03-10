package middlewares

import (
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Sessions() gin.HandlerFunc {
	// Create new cookie store with secure auth
	store := cookie.NewStore([]byte(os.Getenv("AUTH_KEY")))
	store.Options(sessions.Options{
		MaxAge:   60 * 60 * 24 * 30,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return sessions.Sessions("session_user", store)
}
